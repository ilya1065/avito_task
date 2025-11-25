package service

import (
	"avito_task/internal/entity"
	"avito_task/internal/repo"
	"context"
)

type TeamService struct {
	teamRepo repo.Team
}

func NewTeamService(team repo.Team) *TeamService {
	return &TeamService{teamRepo: team}
}

func (s *TeamService) Create(ctx context.Context, team entity.Team, members []entity.User) (entity.Team, []entity.User, error) {
	return s.teamRepo.CreateTeam(ctx, team, members)
}

func (s *TeamService) GetByName(ctx context.Context, name string) (entity.Team, []entity.User, error) {
	return s.teamRepo.GetTeamByName(ctx, name)
}
