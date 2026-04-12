package service

import (
	"errors"
	"grpc-crud/internal/user/model"
	"grpc-crud/internal/user/repo"
)

type UserService interface {
	CreateUser(name, email string) (*model.User, error)
	GetUser(id int64) (*model.User, error)
	UpdateUser(id int64, name, email string) (*model.User, error)
	DeleteUser(id int64) error
}

type userService struct {
	repo repo.UserRepo
}

func NewUserService(r repo.UserRepo) UserService {
	return &userService{repo: r}
}

// CreateUser applies business logic for creating a user record.
func (s *userService) CreateUser(name, email string) (*model.User, error) {
	user := &model.User{Name: name, Email: email}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser retrieves a user from the repository by ID.
func (s *userService) GetUser(id int64) (*model.User, error) {
	return s.repo.GetByID(id)
}

// UpdateUser validates and persists updates to an existing user.
func (s *userService) UpdateUser(id int64, name, email string) (*model.User, error) {
	// Verify the user exists before updating.
	existingUser, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	existingUser.Name = name
	existingUser.Email = email

	if err := s.repo.Update(existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

// DeleteUser removes a user from storage after confirming existence.
func (s *userService) DeleteUser(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}
