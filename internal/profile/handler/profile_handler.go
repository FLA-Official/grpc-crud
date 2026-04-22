package handler

import (
	"context"

	profilev1 "grpc-crud/gen/profile/v1"
	middleware "grpc-crud/internal/middlewares"
	"grpc-crud/internal/profile/model"
	"grpc-crud/internal/profile/service"
	"grpc-crud/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileHandler struct {
	profilev1.UnimplementedProfileServiceServer
	svc *service.ProfileService
}

func NewProfileHandler(s *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		svc: s,
	}
}

func (h *ProfileHandler) CreateProfile(ctx context.Context, req *profilev1.CreateProfileRequest) (*profilev1.ProfileResponse, error) {
	logger := utils.LoggerFromContext(ctx)
	logger.Info("creating profile", "user_id", req.UserId, "name", req.Name)

	// Custom validation example
	if len(req.Name) < 3 {
		logger.Error("profile name validation failed", "user_id", req.UserId, "name_length", len(req.Name))
		return nil, status.Errorf(codes.InvalidArgument, "name too short")
	}

	p := &model.Profile{
		UserID:   req.UserId,
		Name:     req.Name,
		FullName: req.FullName,
		Email:    req.Email,
		Bio:      req.Bio,
	}

	if err := p.Validate(); err != nil {
		logger.Error("profile validation failed", "user_id", req.UserId, "error", err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	if err := h.svc.Create(ctx, p); err != nil {
		logger.Error("failed to create profile", "user_id", req.UserId, "error", err.Error())
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	logger.Info("profile created successfully", "user_id", req.UserId)

	return &profilev1.ProfileResponse{
		Profile: &profilev1.Profile{
			UserId:   p.UserID,
			Name:     p.Name,
			FullName: p.FullName,
			Email:    p.Email,
			Bio:      p.Bio,
		},
	}, nil
}

func (h *ProfileHandler) GetProfile(ctx context.Context, req *profilev1.GetProfileRequest) (*profilev1.ProfileResponse, error) {
	logger := utils.LoggerFromContext(ctx)
	logger.Info("fetching profile", "user_id", req.UserId)

	payload := ctx.Value(middleware.UserContextKey)
	if payload == nil {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}
	userCtx := payload.(*utils.Payload)

	if req.UserId != int64(userCtx.ID) {
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}

	p, err := h.svc.Get(ctx, req.UserId)
	if err != nil {
		logger.Error("failed to fetch profile", "user_id", req.UserId, "error", err.Error())
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}

	logger.Info("profile fetched successfully", "user_id", req.UserId)

	return &profilev1.ProfileResponse{
		Profile: &profilev1.Profile{
			UserId:   p.UserID,
			Name:     p.Name,
			FullName: p.FullName,
			Email:    p.Email,
			Bio:      p.Bio,
		},
	}, nil
}

func (h *ProfileHandler) UpdateProfile(ctx context.Context, req *profilev1.UpdateProfileRequest) (*profilev1.ProfileResponse, error) {
	logger := utils.LoggerFromContext(ctx)
	logger.Info("updating profile", "user_id", req.UserId)

	payload := ctx.Value(middleware.UserContextKey)
	if payload == nil {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}
	userCtx := payload.(*utils.Payload)

	if req.UserId != int64(userCtx.ID) {
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}

	p, err := h.svc.Get(ctx, req.UserId)
	if err != nil {
		logger.Error("failed to retrieve profile for update", "user_id", req.UserId, "error", err.Error())
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}

	p.Name = req.Name
	p.FullName = req.FullName
	p.Bio = req.Bio

	if err := p.Validate(); err != nil {
		logger.Error("profile validation failed during update", "user_id", req.UserId, "error", err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	if err := h.svc.Update(ctx, p); err != nil {
		logger.Error("failed to update profile", "user_id", req.UserId, "error", err.Error())
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	logger.Info("profile updated successfully", "user_id", req.UserId)

	return &profilev1.ProfileResponse{
		Profile: &profilev1.Profile{
			UserId:   p.UserID,
			Name:     p.Name,
			FullName: p.FullName,
			Email:    p.Email,
			Bio:      p.Bio,
		},
	}, nil
}

func (h *ProfileHandler) DeleteProfile(ctx context.Context, req *profilev1.DeleteProfileRequest) (*profilev1.DeleteProfileResponse, error) {
	logger := utils.LoggerFromContext(ctx)
	logger.Info("deleting profile", "user_id", req.UserId)

	payload := ctx.Value(middleware.UserContextKey)
	if payload == nil {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}
	userCtx := payload.(*utils.Payload)

	if req.UserId != int64(userCtx.ID) {
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}

	if err := h.svc.Delete(ctx, req.UserId); err != nil {
		logger.Error("failed to delete profile", "user_id", req.UserId, "error", err.Error())
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	logger.Info("profile deleted successfully", "user_id", req.UserId)
	return &profilev1.DeleteProfileResponse{
		Success: true,
	}, nil
}
