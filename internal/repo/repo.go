package repo

import (
	"avito_task/internal/entity"
	"avito_task/internal/repo/pgdb"
	"context"

	"github.com/jmoiron/sqlx"
)

type User interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (entity.User, error)
	GetReview(ctx context.Context, userID string) (string, []entity.PullRequestShort, error)
}

type Team interface {
	GetTeamByName(ctx context.Context, teamName string) (entity.Team, []entity.User, error)
	CreateTeam(ctx context.Context, team entity.Team, members []entity.User) (entity.Team, []entity.User, error)
}

type PullRequest interface {
	Create(ctx context.Context, pr entity.PullRequestShort) (entity.PullRequestShort, []string, error)
	Merge(ctx context.Context, PrId string) (entity.PullRequest, []string, error)
	Reassign(ctx context.Context, prID string, oldReviewerID string) (entity.PullRequest, []string, string, error)
}

type Repositories struct {
	User
	Team
	PullRequest
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		User:        pgdb.NewUserRepo(db),
		Team:        pgdb.NewTeamRepo(db),
		PullRequest: pgdb.NewPullRequestRepo(db),
	}
}
