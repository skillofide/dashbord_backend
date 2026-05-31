package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	gqlhandler "github.com/graphql-go/handler"
	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pkgauth "github.com/skillofide/pkg/auth"
	pkgdb "github.com/skillofide/pkg/database"
	pkglog "github.com/skillofide/pkg/logger"
	"github.com/skillofide/proto/codec"
	executionv1 "github.com/skillofide/proto/execution/v1"
	problemv1 "github.com/skillofide/proto/problem/v1"
	progressv1 "github.com/skillofide/proto/progress/v1"
	submissionv1 "github.com/skillofide/proto/submission/v1"

	"github.com/skillofide/api-gateway/graph/generated"
	"github.com/skillofide/api-gateway/graph/resolvers"
	"github.com/skillofide/api-gateway/middleware"
)

func main() {
	codec.Register()
	cfg := loadConfig()
	log := pkglog.New(cfg.logLevel)
	defer log.Sync() //nolint:errcheck

	// ── PostgreSQL ──────────────────────────────────────────────────────────
	var pgPool *pgxpool.Pool
	if cfg.postgresDSN != "" {
		pool, err := pkgdb.NewPostgresPool(context.Background(), cfg.postgresDSN)
		if err != nil {
			log.Error("connect postgres for login database", zap.Error(err))
		} else {
			pgPool = pool
			defer pgPool.Close()
			log.Info("postgres connected for users database")

			// Ensure users table and seed data
			if err := ensureUsersTable(context.Background(), pgPool, log); err != nil {
				log.Error("ensure users table failed", zap.Error(err))
			}
		}
	}

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
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

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
	mux.HandleFunc("/login", handleLogin(pgPool, cfg.jwtSecret, log))

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
	notificationServiceURL string
	jwtSecret              string
	jwtPublicKey           string
	postgresDSN            string
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
		notificationServiceURL: env("NOTIFICATION_SERVICE_URL", "http://localhost:8081"),
		jwtSecret:              env("JWT_SECRET", "dev-secret-change-in-production"),
		jwtPublicKey:           env("JWT_PUBLIC_KEY", ""),
		postgresDSN:            env("POSTGRES_DSN", ""),
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

// ─── DB Seeding & Login Handler ──────────────────────────────────────────────

func ensureUsersTable(ctx context.Context, pool *pgxpool.Pool, log *zap.Logger) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email      TEXT NOT NULL UNIQUE,
			name       TEXT NOT NULL,
			password   TEXT NOT NULL,
			role       TEXT NOT NULL DEFAULT 'student',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	// Seed default user
	_, err = pool.Exec(ctx, `
		INSERT INTO users (email, name, password, role)
		VALUES ('admin@skillofied.com', 'Admin User', 'skillofied123', 'admin')
		ON CONFLICT (email) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("seed default user: %w", err)
	}

	log.Info("users table verified and seeded successfully")
	return nil
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

func handleLogin(pool *pgxpool.Pool, jwtSecret string, log *zap.Logger) http.HandlerFunc {
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

		if pool == nil {
			log.Error("login request failed: postgres database unavailable")
			http.Error(w, `{"error":"Database connection unavailable"}`, http.StatusInternalServerError)
			return
		}

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		var user profileInfo
		var dbPassword string
		err := pool.QueryRow(r.Context(), `
			SELECT id::text, email, name, password, role
			FROM   users
			WHERE  email = $1
		`, req.Email).Scan(&user.ID, &user.Email, &user.Name, &dbPassword, &user.Role)

		if err != nil {
			if err == pgx.ErrNoRows {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"Invalid email or password"}`)) //nolint:errcheck
				return
			}
			log.Error("login database query error", zap.Error(err))
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Simple plain text comparison for manual seeds
		if req.Password != dbPassword {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Invalid email or password"}`)) //nolint:errcheck
			return
		}

		// Generate token
		token, err := pkgauth.GenerateToken(user.ID, user.Email, user.Role, jwtSecret, 24*time.Hour)
		if err != nil {
			log.Error("token generation failed", zap.Error(err))
			http.Error(w, `{"error":"Failed to generate auth token"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Token: token,
			User:  user,
		}) //nolint:errcheck
	}
}
