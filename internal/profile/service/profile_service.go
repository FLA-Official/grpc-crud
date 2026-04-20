package service

import (
	"context"
	"errors"
	"grpc-crud/internal/profile/model"
	"grpc-crud/internal/profile/repo"
)

type ProfileService struct {
	repo repo.ProfileRepo
}

func NewProfileService(r repo.ProfileRepo) *ProfileService {
	return &ProfileService{repo: r}
}

func (s *ProfileService) Create(ctx context.Context, p *model.Profile) error {
	if _, err := s.repo.Get(p.UserID); err == nil {
		return errors.New("profile already exists")
	}
	return s.repo.Create(p)
}

func (s *ProfileService) Get(ctx context.Context, userID int64) (*model.Profile, error) {
	p, err := s.repo.Get(userID)
	if err != nil {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (s *ProfileService) Update(ctx context.Context, p *model.Profile) error {
	return s.repo.Update(p)
}

func (s *ProfileService) Delete(ctx context.Context, userID int64) error {
	return s.repo.Delete(userID)
}
