package service

import (
	"context"
	"errors"
	profileModel "grpc-crud/internal/profile/model"
	profileRepo "grpc-crud/internal/profile/repo"
	"grpc-crud/internal/user/model"
	userModel "grpc-crud/internal/user/model"
	"grpc-crud/internal/user/repo"
	"grpc-crud/utils"
)

type UserService interface {
	CreateUser(ctx context.Context, userName, email, password string) (*userModel.User, error)
	GetUser(ctx context.Context, id int64) (*userModel.User, error)
	UpdateUser(ctx context.Context, id int64, userName, email, password string) (*userModel.User, error)
	DeleteUser(ctx context.Context, id int64) error
	Login(ctx context.Context, email, password string) (*model.User, error)
}

type userService struct {
	repo        repo.UserRepo
	profileRepo profileRepo.ProfileRepo
}

func NewUserService(r repo.UserRepo, pr profileRepo.ProfileRepo) UserService {
	return &userService{repo: r, profileRepo: pr}
}

// CreateUser applies business logic for creating a user record.
func (s *userService) CreateUser(ctx context.Context, userName, email, password string) (*userModel.User, error) {
	logger := utils.LoggerFromContext(ctx)

	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Error("failed to hash password", "username", userName, "email", email, "error", err.Error())
		return nil, err
	}

	user := &userModel.User{UserName: userName, Email: email, Password: hashedPassword}
	if err := user.Validate(); err != nil {
		logger.Error("user validation failed", "username", userName, "email", email, "error", err.Error())
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		logger.Error("SERVICE ERROR", "err", err)
		logger.Error("failed to create user in repository", "username", userName, "email", email, "error", err.Error())
		return nil, err
	}

	// Create profile
	profile := &profileModel.Profile{
		UserID:   user.ID,
		Name:     userName,
		FullName: "",
		Email:    email,
		Bio:      "",
	}
	if err := profile.Validate(); err != nil {
		logger.Error("profile validation failed", "user_id", user.ID, "error", err.Error())
		return nil, err
	}
	if err := s.profileRepo.Create(profile); err != nil {
		logger.Error("failed to create profile", "user_id", user.ID, "error", err.Error())
		return nil, err
	}

	logger.Info("user and profile created successfully", "user_id", user.ID, "email", email)
	return user, nil
}

// GetUser retrieves a user from the repository by ID.
func (s *userService) GetUser(ctx context.Context, id int64) (*userModel.User, error) {
	logger := utils.LoggerFromContext(ctx)

	user, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error("failed to retrieve user from repository", "user_id", id, "error", err.Error())
		return nil, err
	}
	return user, nil
}

// UpdateUser validates and persists updates to an existing user.
func (s *userService) UpdateUser(ctx context.Context, id int64, userName, email, password string) (*userModel.User, error) {
	logger := utils.LoggerFromContext(ctx)

	existingUser, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error("failed to retrieve user for update", "user_id", id, "error", err.Error())
		return nil, err
	}
	if existingUser == nil {
		logger.Error("user not found", "user_id", id)
		return nil, errors.New("user not found")
	}

	// Hash the password before updating
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Error("failed to hash password", "user_id", id, "email", email, "error", err.Error())
		return nil, err
	}

	existingUser.UserName = userName
	existingUser.Email = email
	existingUser.Password = hashedPassword

	if err := existingUser.Validate(); err != nil {
		logger.Error("user validation failed during update", "user_id", id, "error", err.Error())
		return nil, err
	}

	if err := s.repo.Update(existingUser); err != nil {
		logger.Error("failed to update user in repository", "user_id", id, "error", err.Error())
		return nil, err
	}

	logger.Info("user updated successfully", "user_id", id, "email", email)
	return existingUser, nil
}

// DeleteUser removes a user from storage after confirming existence.
func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	logger := utils.LoggerFromContext(ctx)

	_, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error("failed to retrieve user for deletion", "user_id", id, "error", err.Error())
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		logger.Error("failed to delete user", "user_id", id, "error", err.Error())
		return err
	}

	logger.Info("user deleted successfully", "user_id", id)
	return nil
}

// login
func (s *userService) Login(ctx context.Context, email, password string) (*model.User, error) {

	logger := utils.LoggerFromContext(ctx)

	user, err := s.repo.Find(email)
	if err != nil {
		logger.Error("user not found", "email", email)
		return nil, errors.New("invalid credentials")
	}

	err = utils.CheckPassword(user.Password, password)
	if err != nil {
		logger.Error("invalid password", "email", email)
		return nil, errors.New("invalid credentials")
	}

	logger.Info("user logged in", "user_id", user.ID)

	return user, nil
}
