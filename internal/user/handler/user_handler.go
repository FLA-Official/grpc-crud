package handler

import (
	"context"
	"fmt"

	userpb "grpc-crud/gen/user/v1"
	"grpc-crud/internal/config"
	middleware "grpc-crud/internal/middlewares"
	"grpc-crud/internal/user/service"
	"grpc-crud/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userService service.UserService
	userpb.UnimplementedUserServiceServer
}

type ReqLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

// Greetings handles the gRPC Greetings request and returns a hello message.
func (h *UserHandler) Greetings(ctx context.Context, req *userpb.HelloRequest) (*userpb.HelloResponse, error) {
	return &userpb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

// CreateUser handles the gRPC CreateUser request and forwards it to the service layer.
func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	logger := utils.LoggerFromContext(ctx)
	logger.Info("creating user", "username", req.UserName, "email", req.Email)

	user, err := h.userService.CreateUser(ctx, req.UserName, req.Email, req.Password)
	if err != nil {
		logger.Error("failed to create user", "username", req.UserName, "email", req.Email, "error", err.Error())
		return nil, err
	}

	logger.Info("user created successfully", "user_id", user.ID, "email", user.Email)
	return &userpb.CreateUserResponse{
		User: &userpb.User{
			Id:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
		},
	}, nil
}

// GetUser handles the gRPC GetUser request and returns a user payload.
func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	logger := utils.LoggerFromContext(ctx)
	logger.Info("fetching user", "user_id", req.Id)

	payload := ctx.Value(middleware.UserContextKey)
	if payload == nil {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}
	userCtx := payload.(*utils.Payload)

	if req.Id != int64(userCtx.ID) {
		return nil, status.Error(codes.PermissionDenied, "you can only access your own data")
	}

	user, err := h.userService.GetUser(ctx, req.Id)
	if err != nil {
		logger.Error("failed to fetch user", "user_id", req.Id, "error", err.Error())
		return nil, err
	}

	logger.Info("user fetched successfully", "user_id", user.ID, "email", user.Email)
	return &userpb.GetUserResponse{
		User: &userpb.User{
			Id:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
		},
	}, nil
}

// UpdateUser handles the gRPC UpdateUser request and returns the updated user.
func (h *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	logger := utils.LoggerFromContext(ctx)
	logger.Info("RAW REQUEST", "req", req)
	logger.Info("updating user", "user_id", req.Id, "email", req.Email)

	payload := ctx.Value(middleware.UserContextKey)
	if payload == nil {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}
	userCtx := payload.(*utils.Payload)

	if req.Id != int64(userCtx.ID) {
		return nil, status.Error(codes.PermissionDenied, "cannot update other users")
	}

	user, err := h.userService.UpdateUser(ctx, req.Id, req.UserName, req.Email, req.Password)
	if err != nil {
		logger.Error("failed to update user", "user_id", req.Id, "error", err.Error())
		return nil, err
	}

	logger.Info("user updated successfully", "user_id", user.ID, "email", user.Email)
	return &userpb.UpdateUserResponse{
		User: &userpb.User{
			Id:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
		},
	}, nil
}

// DeleteUser handles the gRPC DeleteUser request and removes the user record.
func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	logger := utils.LoggerFromContext(ctx)
	logger.Info("deleting user", "user_id", req.Id)

	payload := ctx.Value(middleware.UserContextKey)
	if payload == nil {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}
	userCtx := payload.(*utils.Payload)

	if req.Id != int64(userCtx.ID) {
		return nil, status.Error(codes.PermissionDenied, "cannot delete other users")
	}

	err := h.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		logger.Error("failed to delete user", "user_id", req.Id, "error", err.Error())
		return nil, err
	}

	logger.Info("user deleted successfully", "user_id", req.Id)
	return &userpb.DeleteUserResponse{
		Message: fmt.Sprintf("User with ID %d deleted successfully", req.Id),
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {

	logger := utils.LoggerFromContext(ctx)

	usr, err := h.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.Error("login failed", "email", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	cnf := config.GetConfig()

	accessToken, err := utils.CreateJWT(cnf.JWTSecretKey, utils.Payload{
		ID:       int(usr.ID),
		Username: usr.UserName,
		Email:    usr.Email,
	})
	if err != nil {
		logger.Error("failed to create jwt", "user_id", usr.ID)
		return nil, fmt.Errorf("internal server error")
	}

	logger.Info("user login success", "user_id", usr.ID)

	return &userpb.LoginResponse{
		AccessToken: accessToken,
		User: &userpb.User{
			Id:       usr.ID,
			UserName: usr.UserName,
			Email:    usr.Email,
		},
	}, nil
}
