package hendler

import (
	"avito_task/internal/service"
	"net/http"
)

type User interface {
	SetIsActive(w http.ResponseWriter, r *http.Request)
	GetReview(w http.ResponseWriter, r *http.Request)
}
type Team interface {
	CreateTeam(w http.ResponseWriter, r *http.Request)
	GetTeam(w http.ResponseWriter, r *http.Request)
}

type PullRequest interface {
	Create(w http.ResponseWriter, r *http.Request)
	Merge(w http.ResponseWriter, r *http.Request)
	Reassign(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	UserHandler User
	TeamHandler Team
	PrHandler   PullRequest
}

func NewHendler(service *service.Service) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(service.UserService),
		TeamHandler: NewTeamHendler(service.TeamService),
		PrHandler:   NewPullRequestHendler(service.PullRequestService),
	}
}
