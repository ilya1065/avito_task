package service

import (
	"avito_task/internal/entity"
	"avito_task/internal/repo"
	"context"
)

type UserService struct {
	userRepo repo.User
}

func NewUserService(userRepo repo.User) *UserService { return &UserService{userRepo: userRepo} }

func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (entity.User, error) {
	return s.userRepo.SetIsActive(ctx, userID, isActive)
}

func (s *UserService) GetReview(ctx context.Context, userID string) (string, []entity.PullRequestShort, error) {
	return s.userRepo.GetReview(ctx, userID)
}
