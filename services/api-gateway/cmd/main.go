package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	gqlhandler "github.com/graphql-go/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pkgauth "github.com/skillofide/pkg/auth"
	pkglog "github.com/skillofide/pkg/logger"
	"github.com/skillofide/proto/codec"
	executionv1 "github.com/skillofide/proto/execution/v1"
	problemv1 "github.com/skillofide/proto/problem/v1"
	progressv1 "github.com/skillofide/proto/progress/v1"
	submissionv1 "github.com/skillofide/proto/submission/v1"
	userv1 "github.com/skillofide/proto/user/v1"

	"github.com/skillofide/api-gateway/graph/generated"
	"github.com/skillofide/api-gateway/graph/resolvers"
	"github.com/skillofide/api-gateway/middleware"
)

func main() {
	codec.Register()
	cfg := loadConfig()
	log := pkglog.New(cfg.logLevel)
	defer log.Sync() //nolint:errcheck

	// ── User Service gRPC Connection ──────────────────────────────────────────
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	userConn := mustDial(cfg.userServiceAddr, dialOpts, log)
	defer userConn.Close()
	userSvcClient := userv1.NewUserServiceClient(userConn)

	// ── Auth validator ────────────────────────────────────────────────────────
	var jwtValidator *pkgauth.Validator
	if cfg.jwtPublicKey != "" {
		v, err := pkgauth.NewRS256Validator(cfg.jwtPublicKey)
		if err != nil {
			log.Fatal("init RS256 validator", zap.Error(err))
		}
		jwtValidator = v
	} else {
		jwtValidator = pkgauth.NewHMACValidator(cfg.jwtSecret)
	}

	// ── gRPC connections ──────────────────────────────────────────────────────

	probConn := mustDial(cfg.problemServiceAddr, dialOpts, log)
	defer probConn.Close()

	subConn := mustDial(cfg.submissionServiceAddr, dialOpts, log)
	defer subConn.Close()

	execConn := mustDial(cfg.executionServiceAddr, dialOpts, log)
	defer execConn.Close()

	progConn := mustDial(cfg.progressServiceAddr, dialOpts, log)
	defer progConn.Close()

	// ── Build GraphQL schema ──────────────────────────────────────────────────
	clients := &generated.Clients{
		Problems: &resolvers.ProblemClients{
			ProblemSvc: problemv1.NewProblemServiceClient(probConn),
			Log:        log,
		},
		Submissions: &resolvers.SubmissionClients{
			SubmissionSvc: submissionv1.NewSubmissionServiceClient(subConn),
			ExecutionSvc:  executionv1.NewExecutionServiceClient(execConn),
			Log:           log,
		},
		Progress: &resolvers.ProgressClients{
			ProgressSvc: progressv1.NewProgressServiceClient(progConn),
			Log:         log,
		},
		User: &resolvers.UserClients{
			UserSvc: userSvcClient,
			Log:     log,
		},
	}

	schema, err := generated.BuildSchema(clients)
	if err != nil {
		log.Fatal("build graphql schema", zap.Error(err))
	}

	// ── HTTP routes ───────────────────────────────────────────────────────────
	gqlHandler := gqlhandler.New(&gqlhandler.Config{
		Schema:     &schema,
		Pretty:     cfg.devMode,
		GraphiQL:   cfg.devMode, // enable playground in dev
		Playground: cfg.devMode,
	})

	mux := http.NewServeMux()

	// GraphQL endpoint
	mux.Handle("/graphql", gqlHandler)

	// REST Login endpoint
	mux.HandleFunc("/login", handleLogin(userSvcClient, cfg.jwtSecret, log))

	// REST Profile endpoints
	mux.HandleFunc("/profile", handleProfile(userSvcClient, log))

	// WebSocket proxy to notification-service
	if cfg.notificationServiceURL != "" {
		wsTarget, err := url.Parse(cfg.notificationServiceURL)
		if err == nil {
			proxy := httputil.NewSingleHostReverseProxy(wsTarget)
			mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
				proxy.ServeHTTP(w, r)
			})
		}
	}

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`)) //nolint:errcheck
	})

	// ── Apply middleware chain ────────────────────────────────────────────────
	authMW := middleware.Auth(jwtValidator, log, "/health", "/login") // JWT is parsed for /graphql; resolvers decide per-field auth
	corsMW := middleware.CORS(cfg.allowedOrigins)

	handler := corsMW(authMW(mux))

	server := &http.Server{
		Addr:         ":" + cfg.httpPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Info("api-gateway starting",
		zap.String("port", cfg.httpPort),
		zap.Bool("graphiql", cfg.devMode),
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server error", zap.Error(err))
		}
	}()

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down api-gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx) //nolint:errcheck
	log.Info("api-gateway stopped")
}

func mustDial(addr string, opts []grpc.DialOption, log *zap.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatal("dial gRPC service", zap.String("addr", addr), zap.Error(err))
	}
	log.Info("gRPC connection established", zap.String("addr", addr))
	return conn
}

type config struct {
	httpPort               string
	problemServiceAddr     string
	submissionServiceAddr  string
	executionServiceAddr   string
	progressServiceAddr    string
	userServiceAddr        string
	notificationServiceURL string
	jwtSecret              string
	jwtPublicKey           string
	allowedOrigins         string
	devMode                bool
	logLevel               string
}

func loadConfig() config {
	return config{
		httpPort:               env("HTTP_PORT", "8080"),
		problemServiceAddr:     env("PROBLEM_SERVICE_ADDR", "localhost:50051"),
		submissionServiceAddr:  env("SUBMISSION_SERVICE_ADDR", "localhost:50053"),
		executionServiceAddr:   env("EXECUTION_SERVICE_ADDR", "localhost:50052"),
		progressServiceAddr:    env("PROGRESS_SERVICE_ADDR", "localhost:50054"),
		userServiceAddr:        env("USER_SERVICE_ADDR", "localhost:50055"),
		notificationServiceURL: env("NOTIFICATION_SERVICE_URL", "http://localhost:8081"),
		jwtSecret:              env("JWT_SECRET", "dev-secret-change-in-production"),
		jwtPublicKey:           env("JWT_PUBLIC_KEY", ""),
		allowedOrigins:         env("ALLOWED_ORIGINS", "*"),
		devMode:                env("DEV_MODE", "true") == "true",
		logLevel:               env("LOG_LEVEL", "info"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ─── Login Handler ──────────────────────────────────────────────

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ─── Profile Handler ─────────────────────────────────────────────

func handleProfile(userSvc userv1.UserServiceClient, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// OPTIONS (CORS preflight handled by CORS middleware, but handle gracefully)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Authenticated user
		userID := middleware.UserIDFromContext(r.Context())
		if userID == "" {
			http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodGet:
			resp, err := userSvc.GetProfile(r.Context(), &userv1.GetProfileRequest{UserID: userID})
			if err != nil {
				log.Error("get profile failed", zap.Error(err))
				http.Error(w, `{"error":"failed to get profile"}`, http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(resp.Profile) //nolint:errcheck

		case http.MethodPut:
			var p userv1.UserProfile
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
				return
			}
			p.UserID = userID // enforce user_id from JWT, never from body
			resp, err := userSvc.UpsertProfile(r.Context(), &userv1.UpsertProfileRequest{Profile: &p})
			if err != nil {
				log.Error("upsert profile failed", zap.Error(err))
				http.Error(w, `{"error":"failed to save profile"}`, http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(resp) //nolint:errcheck

		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

type loginResponse struct {
	Token string      `json:"token"`
	User  profileInfo `json:"user"`
}

type profileInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func handleLogin(userSvc userv1.UserServiceClient, jwtSecret string, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle CORS preflight options
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		resp, err := userSvc.VerifyUser(r.Context(), &userv1.VerifyUserRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"Invalid email or password"}`)) //nolint:errcheck
				return
			}
			log.Error("login verification failed", zap.Error(err))
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Generate token
		token, err := pkgauth.GenerateToken(resp.Id, resp.Email, resp.Role, jwtSecret, 24*time.Hour)
		if err != nil {
			log.Error("token generation failed", zap.Error(err))
			http.Error(w, `{"error":"Failed to generate auth token"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Token: token,
			User: profileInfo{
				ID:    resp.Id,
				Email: resp.Email,
				Name:  resp.Name,
				Role:  resp.Role,
			},
		}) //nolint:errcheck
	}
}
