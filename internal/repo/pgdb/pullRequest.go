package pgdb

import (
	"avito_task/internal/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type PullRequestRepo struct {
	db *sqlx.DB
}

func NewPullRequestRepo(db *sqlx.DB) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

func (r *PullRequestRepo) Create(ctx context.Context, pr entity.PullRequestShort) (entity.PullRequestShort, []string, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return entity.PullRequestShort{}, nil, err
	}
	defer tx.Rollback()

	// Проверяем, что PR с таким ID ещё не существует
	var exists bool
	err = tx.GetContext(ctx, &exists,
		`SELECT EXISTS (
           SELECT 1 FROM pull_requests WHERE pull_request_id = $1
        )`,
		pr.PullRequestID,
	)
	if err != nil {
		return entity.PullRequestShort{}, nil, err
	}
	if exists {
		return entity.PullRequestShort{}, nil, errors.New("PR_EXISTS")
	}

	// Получаем команду автора
	var teamName string
	err = tx.GetContext(ctx, &teamName,
		`SELECT team_name FROM users WHERE user_id = $1`,
		pr.AuthorID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// автор не найден
			return entity.PullRequestShort{}, nil, errors.New("NOT_FOUND")
		}
		return entity.PullRequestShort{}, nil, err
	}

	// Выбираем до двух активных ревьюверов из команды автора, исключая самого автора
	var reviewerIDs []string
	err = tx.SelectContext(ctx, &reviewerIDs,
		`SELECT user_id
         FROM users
         WHERE team_name = $1
           AND is_active = TRUE
           AND user_id <> $2
         ORDER BY random()
         LIMIT 2`,
		teamName,
		pr.AuthorID,
	)
	if err != nil {
		return entity.PullRequestShort{}, nil, err
	}

	// Ставим статус OPEN
	pr.Status = entity.PROpen

	// Вставляем PR
	_, err = tx.ExecContext(ctx,
		`INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
         VALUES ($1, $2, $3, $4)`,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
	)
	if err != nil {
		return entity.PullRequestShort{}, nil, err
	}

	// Вставляем ревьюверов в pr_reviewers
	for _, rID := range reviewerIDs {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
             VALUES ($1, $2)`,
			pr.PullRequestID,
			rID,
		)
		if err != nil {
			return entity.PullRequestShort{}, nil, err
		}
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return entity.PullRequestShort{}, nil, err
	}

	return pr, reviewerIDs, nil
}

func (r *PullRequestRepo) Merge(ctx context.Context, PrId string) (entity.PullRequest, []string, error) {
	var status string
	err := r.db.GetContext(ctx, &status, `SELECT status FROM pull_requests WHERE pull_request_id = $1`, PrId)
	if err != nil {
		return entity.PullRequest{}, nil, fmt.Errorf("NOT_FOUND")
	}
	switch status {
	case "MERGED":
		{
			var assignedReviewers []string
			var pr entity.PullRequest
			err = r.db.GetContext(ctx, &pr, `SELECT * FROM Pull_requests WHERE pull_request_id = $1`, PrId)
			if err != nil {
				return entity.PullRequest{}, nil, fmt.Errorf("NOT_FOUND")
			}
			err = r.db.SelectContext(ctx, &assignedReviewers, `SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1`, PrId)
			if err != nil {
				return entity.PullRequest{}, nil, fmt.Errorf("NOT_FOUND")
			}

			return pr, assignedReviewers, err
		}
	case "OPEN":
		{
			st := entity.PRMerged
			var pr entity.PullRequest
			var assignedReviewers []string
			_, err = r.db.ExecContext(ctx, `UPDATE pull_requests SET status = $1,merged_at =$2  WHERE pull_request_id = $3`, st, time.Now(), PrId)
			if err != nil {
				return entity.PullRequest{}, nil, err
			}
			err = r.db.SelectContext(ctx, &assignedReviewers, `SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1`, PrId)
			err = r.db.GetContext(ctx, &pr, `SELECT * FROM pull_requests WHERE pull_request_id = $1`, PrId)
			if err != nil {
				return entity.PullRequest{}, nil, err
			}
			return pr, assignedReviewers, nil
		}
	default:
		return entity.PullRequest{}, nil, fmt.Errorf("PR_MERGED")

	}

}

func (r *PullRequestRepo) Reassign(ctx context.Context, prID string, oldReviewerID string) (entity.PullRequest, []string, string, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return entity.PullRequest{}, nil, "", err
	}
	defer tx.Rollback()

	// Проверяем PR
	var pr entity.PullRequest
	err = tx.GetContext(ctx, &pr,
		`SELECT pull_request_id,
                pull_request_name,
                author_id,
                status,
                created_at,
                merged_at
         FROM pull_requests
         WHERE pull_request_id = $1`,
		prID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.PullRequest{}, nil, "", errors.New("NOT_FOUND")
		}
		return entity.PullRequest{}, nil, "", err
	}

	// проверяем PR уже MERGED?
	if pr.Status == entity.PRMerged {
		return entity.PullRequest{}, nil, "", errors.New("PR_MERGED")
	}

	// Проверяем что пользователь вообще существует
	var userExists bool
	err = tx.GetContext(ctx, &userExists,
		`SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`,
		oldReviewerID,
	)
	if err != nil {
		return entity.PullRequest{}, nil, "", err
	}
	if !userExists {
		return entity.PullRequest{}, nil, "", errors.New("NOT_FOUND")
	}

	// Проверяем, что он действительно назначен ревьювером этого PR
	var assigned bool
	err = tx.GetContext(ctx, &assigned,
		`SELECT EXISTS(
             SELECT 1
             FROM pr_reviewers
             WHERE pull_request_id = $1 AND reviewer_id = $2
         )`,
		prID, oldReviewerID,
	)
	if err != nil {
		return entity.PullRequest{}, nil, "", err
	}
	if !assigned {
		return entity.PullRequest{}, nil, "", errors.New("NOT_ASSIGNED")
	}

	// Узнаём команду старого ревьювера
	var teamName string
	err = tx.GetContext(ctx, &teamName,
		`SELECT team_name FROM users WHERE user_id = $1`,
		oldReviewerID,
	)
	if err != nil {
		return entity.PullRequest{}, nil, "", err
	}

	// Ищем нового ревьювера в его команде
	var newReviewerID string
	err = tx.GetContext(ctx, &newReviewerID,
		`SELECT user_id
         FROM users
         WHERE team_name = $1
           AND is_active = TRUE
           AND user_id <> $2
           AND user_id <> $3
           AND user_id NOT IN (
               SELECT reviewer_id
               FROM pr_reviewers
               WHERE pull_request_id = $4
           )
         ORDER BY random()
         LIMIT 1`,
		teamName,
		oldReviewerID,
		pr.AuthorID,
		prID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.PullRequest{}, nil, "", errors.New("NO_CANDIDATE")
		}
		return entity.PullRequest{}, nil, "", err
	}

	// Удаляем старого ревьювера
	_, err = tx.ExecContext(ctx,
		`DELETE FROM pr_reviewers
         WHERE pull_request_id = $1 AND reviewer_id = $2`,
		prID, oldReviewerID,
	)
	if err != nil {
		return entity.PullRequest{}, nil, "", err
	}

	// Добавляем нового ревьювера
	_, err = tx.ExecContext(ctx,
		`INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
         VALUES ($1, $2)`,
		prID, newReviewerID,
	)
	if err != nil {
		return entity.PullRequest{}, nil, "", err
	}

	// Получаем актуальный список ревьюверов
	var reviewers []string
	err = tx.SelectContext(ctx, &reviewers,
		`SELECT reviewer_id
         FROM pr_reviewers
         WHERE pull_request_id = $1`,
		prID,
	)
	if err != nil {
		return entity.PullRequest{}, nil, "", err
	}

	if err = tx.Commit(); err != nil {
		return entity.PullRequest{}, nil, "", err
	}

	return pr, reviewers, newReviewerID, nil
}
