package pgdb

import (
	"avito_task/internal/entity"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) SetIsActive(ctx context.Context, userID string, isActive bool) (entity.User, error) {
	res, err := r.db.ExecContext(ctx, `update users set is_active = $1 where user_id = $2`, isActive, userID)
	if err != nil {
		return entity.User{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return entity.User{}, err
	}
	if rows == 0 {
		return entity.User{}, fmt.Errorf("NOT_FOUND")
	}
	var user entity.User
	err = r.db.GetContext(ctx, &user, `SELECT user_id, username, team_name, is_active FROM users where user_id = $1`, userID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepo) GetReview(ctx context.Context, userID string) (string, []entity.PullRequestShort, error) {
	// Проверка что пользователь существует
	var exists bool
	err := r.db.GetContext(ctx, &exists,
		`SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`,
		userID)
	if err != nil {
		return "", nil, err
	}
	if !exists {
		return "", nil, fmt.Errorf("NOT_FOUND")
	}

	// Получение PR, где пользователь является ревьюером
	var prs []entity.PullRequestShort
	err = r.db.SelectContext(ctx, &prs,
		`SELECT 
            pr.pull_request_id,
            pr.pull_request_name,
            pr.author_id,
            pr.status
         FROM pull_requests pr
         JOIN pr_reviewers r 
              ON r.pull_request_id = pr.pull_request_id
         WHERE r.reviewer_id = $1`,
		userID)

	if err != nil {
		return userID, nil, err
	}

	return userID, prs, nil
}
