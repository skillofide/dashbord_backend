package handler

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/skillofide/user-service/internal/repository"
	userv1 "github.com/skillofide/proto/user/v1"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	repo   *repository.UserRepository
	logger *zap.Logger
}

func New(repo *repository.UserRepository, log *zap.Logger) *UserHandler {
	return &UserHandler{repo: repo, logger: log}
}

func (h *UserHandler) VerifyUser(ctx context.Context, req *userv1.VerifyUserRequest) (*userv1.VerifyUserResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	resp, err := h.repo.VerifyUser(ctx, req.Email, req.Password)
	if err != nil {
		if err == pgx.ErrNoRows || err.Error() == "invalid password" {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		h.logger.Error("verify user query failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "verify user failed: %v", err)
	}

	return resp, nil
}

func (h *UserHandler) CreateOrUpdateUser(ctx context.Context, req *userv1.CreateOrUpdateUserRequest) (*userv1.CreateOrUpdateUserResponse, error) {
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email, name and password are required")
	}

	role := req.Role
	if role == "" {
		role = "student"
	}

	err := h.repo.CreateOrUpdateUser(ctx, req.Email, req.Name, req.Password, role)
	if err != nil {
		h.logger.Error("upsert user failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "upsert user failed: %v", err)
	}

	return &userv1.CreateOrUpdateUserResponse{
		Success: true,
		Message: "User upserted successfully",
	}, nil
}

// GetProfile returns the profile for the given user_id.
// If no profile row exists yet an empty profile struct is returned (not an error).
func (h *UserHandler) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.GetProfileResponse, error) {
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	profile, err := h.repo.GetProfile(ctx, req.UserID)
	if err != nil {
		h.logger.Error("get profile failed", zap.String("user_id", req.UserID), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "get profile failed: %v", err)
	}

	if profile == nil {
		// Return an empty profile rather than an error so the frontend can render the empty state
		profile = &userv1.UserProfile{UserID: req.UserID}
	}

	return &userv1.GetProfileResponse{Profile: profile}, nil
}

// UpsertProfile inserts or updates a user's profile.
func (h *UserHandler) UpsertProfile(ctx context.Context, req *userv1.UpsertProfileRequest) (*userv1.UpsertProfileResponse, error) {
	if req.Profile == nil || req.Profile.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "profile with a valid user_id is required")
	}

	if err := h.repo.UpsertProfile(ctx, req.Profile); err != nil {
		h.logger.Error("upsert profile failed", zap.String("user_id", req.Profile.UserID), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "upsert profile failed: %v", err)
	}

	return &userv1.UpsertProfileResponse{
		Success: true,
		Message: "Profile saved successfully",
	}, nil
}

