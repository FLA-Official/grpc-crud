package handler

import (
	"context"

	profilev1 "grpc-crud/gen/profile/v1"
	"grpc-crud/internal/profile/model"
	"grpc-crud/internal/profile/service"

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

	// Custom validation example
	if len(req.Name) < 3 {
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
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	if err := h.svc.Create(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

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
	p, err := h.svc.Get(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}

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
	p, err := h.svc.Get(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}

	p.Name = req.Name
	p.FullName = req.FullName
	p.Bio = req.Bio

	if err := p.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	if err := h.svc.Update(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

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
	if err := h.svc.Delete(ctx, req.UserId); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &profilev1.DeleteProfileResponse{
		Success: true,
	}, nil
}
