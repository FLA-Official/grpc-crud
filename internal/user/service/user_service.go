package service

import (
	"errors"
	profileModel "grpc-crud/internal/profile/model"
	profileRepo "grpc-crud/internal/profile/repo"
	userModel "grpc-crud/internal/user/model"
	"grpc-crud/internal/user/repo"
)

type UserService interface {
	CreateUser(userName, email, password, fullName string) (*userModel.User, error)
	GetUser(id int64) (*userModel.User, error)
	UpdateUser(id int64, userName, email, password, fullName string) (*userModel.User, error)
	DeleteUser(id int64) error
}

type userService struct {
	repo        repo.UserRepo
	profileRepo profileRepo.ProfileRepo
}

func NewUserService(r repo.UserRepo, pr profileRepo.ProfileRepo) UserService {
	return &userService{repo: r, profileRepo: pr}
}

// CreateUser applies business logic for creating a user record.
func (s *userService) CreateUser(userName, email, password, fullName string) (*userModel.User, error) {
	user := &userModel.User{UserName: userName, Email: email, Password: password}
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Create profile
	profile := &profileModel.Profile{
		UserID:   user.ID,
		Name:     userName,
		FullName: fullName,
		Email:    email,
		Bio:      "",
	}
	if err := profile.Validate(); err != nil {
		return nil, err
	}
	if err := s.profileRepo.Create(profile); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user from the repository by ID.
func (s *userService) GetUser(id int64) (*userModel.User, error) {
	return s.repo.GetByID(id)
}

// UpdateUser validates and persists updates to an existing user.
func (s *userService) UpdateUser(id int64, userName, email, password, fullName string) (*userModel.User, error) {
	existingUser, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	existingUser.UserName = userName
	existingUser.Email = email
	existingUser.Password = password

	if err := existingUser.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(existingUser); err != nil {
		return nil, err
	}

	// Update profile if fullName provided
	if fullName != "" {
		profile, err := s.profileRepo.Get(existingUser.ID)
		if err != nil {
			return nil, err
		}
		if profile != nil {
			profile.Name = userName
			profile.FullName = fullName
			if err := profile.Validate(); err != nil {
				return nil, err
			}
			if err := s.profileRepo.Update(profile); err != nil {
				return nil, err
			}
		}
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
