package service

import (
	"avito_task/internal/entity"
	"avito_task/internal/repo"
	"context"
)

type PullRequestService struct {
	pullRequestRepo repo.PullRequest
}

func NewPullRequestService(pullRequestRepo repo.PullRequest) *PullRequestService {
	return &PullRequestService{pullRequestRepo: pullRequestRepo}
}

func (s *PullRequestService) Create(ctx context.Context, pr entity.PullRequestShort) (entity.PullRequestShort, []string, error) {
	return s.pullRequestRepo.Create(ctx, pr)
}
func (s *PullRequestService) Merge(ctx context.Context, id string) (entity.PullRequest, []string, error) {
	return s.pullRequestRepo.Merge(ctx, id)
}
func (s *PullRequestService) Reassign(ctx context.Context, prID, oldReviewerID string) (entity.PullRequest, []string, string, error) {
	return s.pullRequestRepo.Reassign(ctx, prID, oldReviewerID)
}
