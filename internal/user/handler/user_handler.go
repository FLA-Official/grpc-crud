package handler

import (
	"context"
	"fmt"

	"grpc-crud/internal/user/service"
	"grpc-crud/pkg/pb/userpb"
)

type UserHandler struct {
	userService service.UserService
	userpb.UnimplementedUserServiceServer
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

// CreateUser handles the gRPC CreateUser request and forwards it to the service layer.
func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	user, err := h.userService.CreateUser(req.Name, req.Email)
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		User: &userpb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

// GetUser handles the gRPC GetUser request and returns a user payload.
func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := h.userService.GetUser(req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserResponse{
		User: &userpb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

// UpdateUser handles the gRPC UpdateUser request and returns the updated user.
func (h *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	user, err := h.userService.UpdateUser(req.Id, req.Name, req.Email)
	if err != nil {
		return nil, err
	}

	return &userpb.UpdateUserResponse{
		User: &userpb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

// DeleteUser handles the gRPC DeleteUser request and removes the user record.
func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	err := h.userService.DeleteUser(req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.DeleteUserResponse{
		Message: fmt.Sprintf("User with ID %d deleted successfully", req.Id),
	}, nil
}
