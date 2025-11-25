package pgdb

import (
	"avito_task/internal/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TeamRepo struct {
	db *sqlx.DB
}

func NewTeamRepo(db *sqlx.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) CreateTeam(
	ctx context.Context,
	team entity.Team,
	members []entity.User,
) (entity.Team, []entity.User, error) {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return entity.Team{}, nil, err
	}
	defer tx.Rollback()

	//  Проверяем существует ли команда
	var exists bool
	err = tx.GetContext(ctx, &exists,
		`SELECT true FROM teams WHERE team_name=$1`,
		team.TeamName,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return entity.Team{}, nil, err
	}

	if exists {
		return entity.Team{}, nil, fmt.Errorf("TEAM_EXISTS")
	}

	//  Создаём команду
	_, err = tx.ExecContext(ctx,
		`INSERT INTO teams (team_name)
         VALUES ($1)`,
		team.TeamName,
	)
	if err != nil {
		return entity.Team{}, nil, err
	}

	// Создаём/обновляем пользователей
	for i := range members {

		members[i].TeamName = team.TeamName

		_, err = tx.ExecContext(ctx,
			`INSERT INTO users (user_id, username, team_name, is_active)
             VALUES ($1, $2, $3, $4)
             ON CONFLICT (user_id) DO UPDATE
             SET username  = EXCLUDED.username,
                 team_name = EXCLUDED.team_name,
                 is_active = EXCLUDED.is_active`,
			members[i].UserID,
			members[i].Username,
			members[i].TeamName,
			members[i].IsActive,
		)
		if err != nil {
			return entity.Team{}, nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return entity.Team{}, nil, err
	}

	return team, members, nil
}

func (r *TeamRepo) GetTeamByName(
	ctx context.Context,
	teamName string,
) (entity.Team, []entity.User, error) {

	var team entity.Team

	err := r.db.GetContext(ctx, &team,
		`SELECT team_name FROM teams WHERE team_name = $1`,
		teamName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Team{}, nil, fmt.Errorf("NOT_FOUND")
		}
		return entity.Team{}, nil, err
	}

	var members []entity.User
	err = r.db.SelectContext(ctx, &members,
		`SELECT user_id, username, team_name, is_active
         FROM users
         WHERE team_name = $1`,
		teamName,
	)
	if err != nil {
		return entity.Team{}, nil, err
	}

	return team, members, nil
}
