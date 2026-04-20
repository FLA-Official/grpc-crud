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
		UserID: req.UserId,
		Name:   req.Name,
		Email:  req.Email,
		Bio:    req.Bio,
	}

	if err := h.svc.Create(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &profilev1.ProfileResponse{
		Profile: &profilev1.Profile{
			UserId: p.UserID,
			Name:   p.Name,
			Email:  p.Email,
			Bio:    p.Bio,
		},
	}, nil
}
