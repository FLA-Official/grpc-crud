package service

import (
	"context"
	"errors"
	"grpc-crud/internal/profile/model"
	"grpc-crud/internal/profile/repo"
	"grpc-crud/utils"
)

type ProfileService struct {
	repo repo.ProfileRepo
}

func NewProfileService(r repo.ProfileRepo) *ProfileService {
	return &ProfileService{repo: r}
}

func (s *ProfileService) Create(ctx context.Context, p *model.Profile) error {
	logger := utils.LoggerFromContext(ctx)
	
	if _, err := s.repo.Get(p.UserID); err == nil {
		logger.Error("profile already exists", "user_id", p.UserID)
		return errors.New("profile already exists")
	}
	
	if err := s.repo.Create(p); err != nil {
		logger.Error("failed to create profile in repository", "user_id", p.UserID, "error", err.Error())
		return err
	}
	
	logger.Info("profile created successfully", "user_id", p.UserID)
	return nil
}

func (s *ProfileService) Get(ctx context.Context, userID int64) (*model.Profile, error) {
	logger := utils.LoggerFromContext(ctx)
	
	p, err := s.repo.Get(userID)
	if err != nil {
		logger.Error("profile not found", "user_id", userID, "error", err.Error())
		return nil, errors.New("not found")
	}
	return p, nil
}

func (s *ProfileService) Update(ctx context.Context, p *model.Profile) error {
	logger := utils.LoggerFromContext(ctx)
	
	if err := s.repo.Update(p); err != nil {
		logger.Error("failed to update profile in repository", "user_id", p.UserID, "error", err.Error())
		return err
	}
	
	logger.Info("profile updated successfully", "user_id", p.UserID)
	return nil
}

func (s *ProfileService) Delete(ctx context.Context, userID int64) error {
	logger := utils.LoggerFromContext(ctx)
	
	if err := s.repo.Delete(userID); err != nil {
		logger.Error("failed to delete profile", "user_id", userID, "error", err.Error())
		return err
	}
	
	logger.Info("profile deleted successfully", "user_id", userID)
	return nil
}
