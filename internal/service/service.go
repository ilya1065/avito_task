package service

import (
	"avito_task/internal/entity"
	"avito_task/internal/repo"
	"context"
)

type Team interface {
	Create(ctx context.Context, team entity.Team, members []entity.User) (entity.Team, []entity.User, error)
	GetByName(ctx context.Context, name string) (entity.Team, []entity.User, error)
}

type User interface {
	GetReview(ctx context.Context, userID string) (string, []entity.PullRequestShort, error)
	SetIsActive(ctx context.Context, userID string, isActive bool) (entity.User, error)
}

type PullRequest interface {
	Create(ctx context.Context, pr entity.PullRequestShort) (entity.PullRequestShort, []string, error)
	Merge(ctx context.Context, PrId string) (entity.PullRequest, []string, error)
	Reassign(ctx context.Context, prID, oldReviewerID string) (entity.PullRequest, []string, string, error)
}

type Service struct {
	TeamService        Team
	UserService        User
	PullRequestService PullRequest
}

func NewService(repos *repo.Repositories) *Service {
	return &Service{
		TeamService:        NewTeamService(repos.Team),
		UserService:        NewUserService(repos.User),
		PullRequestService: NewPullRequestService(repos.PullRequest),
	}
}
