package hendler

import (
	"avito_task/internal/controller/http/v1/hendler/hendlererrors"
	"avito_task/internal/entity"
	"avito_task/internal/service"
	"encoding/json"
	"net/http"
)

type TeamHendler struct {
	TeamSvc service.Team
}

func NewTeamHendler(teamSvc service.Team) *TeamHendler {
	return &TeamHendler{TeamSvc: teamSvc}
}

type teamMemberDTO struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type creteTeamRequest struct {
	TeamName string          `json:"team_name"`
	Members  []teamMemberDTO `json:"members"`
}

type teamResponse struct {
	Team struct {
		TeamName string          `json:"team_name"`
		Members  []teamMemberDTO `json:"members"`
	} `json:"team"`
}

func (h *TeamHendler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var req creteTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hendlererrors.WriteError(w, http.StatusBadRequest, "INTERNAL_ERROR", err.Error())
		return
	}
	members := make([]entity.User, len(req.Members))
	for i, m := range req.Members {
		members[i] = entity.User{
			UserID:   m.UserID,
			Username: m.UserName,
			IsActive: m.IsActive,
			TeamName: req.TeamName,
		}
	}
	team := entity.Team{TeamName: req.TeamName}
	creteTeam, createMembers, err := h.TeamSvc.Create(r.Context(), team, members)
	if err != nil {
		if err.Error() == "TEAM_EXISTS" {
			hendlererrors.WriteError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
			return
		}
		hendlererrors.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "error creating team")
		return
	}
	resp := teamResponse{}
	resp.Team.TeamName = creteTeam.TeamName
	resp.Team.Members = make([]teamMemberDTO, len(createMembers))

	for i, m := range createMembers {
		resp.Team.Members[i] = teamMemberDTO{
			UserID:   m.UserID,
			UserName: m.Username,
			IsActive: m.IsActive,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *TeamHendler) GetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	TeamName := r.URL.Query().Get("team_name")
	if TeamName == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var team entity.Team
	var members []entity.User
	team, members, err := h.TeamSvc.GetByName(r.Context(), TeamName)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			hendlererrors.WriteError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var resp teamResponse
	resp.Team.TeamName = team.TeamName
	resp.Team.Members = make([]teamMemberDTO, len(members))
	for i, m := range members {
		resp.Team.Members[i] = teamMemberDTO{
			UserID:   m.UserID,
			UserName: m.Username,
			IsActive: m.IsActive,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp.Team)

}
