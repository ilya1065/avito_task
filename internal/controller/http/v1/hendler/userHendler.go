package hendler

import (
	"avito_task/internal/controller/http/v1/hendler/hendlererrors"
	"avito_task/internal/entity"
	"avito_task/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	userSvc service.User
}

func NewUserHandler(service service.User) *UserHandler {
	return &UserHandler{userSvc: service}
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
type SetIsActiveResponse struct {
	User entity.User `json:"user"`
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var req SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fmt.Println(req)
	user, err := h.userSvc.SetIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			hendlererrors.WriteError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var resp SetIsActiveResponse
	resp.User = user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

type PullRequestsDTO struct {
	PullRequestsID   string          `json:"pull_request_id"`
	PullRequestsName string          `json:"pull_request_name"`
	AuthorID         string          `json:"author_id"`
	Status           entity.PRStatus `json:"status"`
}
type GetreviewResponse struct {
	UserID       string            `json:"user_id"`
	PullRequests []PullRequestsDTO `json:"pull_requests"`
}

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var PRShort []entity.PullRequestShort
	userID := r.URL.Query().Get("user_id")
	userID, PRShort, err := h.userSvc.GetReview(r.Context(), userID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			hendlererrors.WriteError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var resp GetreviewResponse
	resp.UserID = userID
	resp.PullRequests = make([]PullRequestsDTO, len(PRShort))
	for i, v := range PRShort {
		resp.PullRequests[i] = PullRequestsDTO{
			PullRequestsID:   v.PullRequestID,
			PullRequestsName: v.PullRequestName,
			AuthorID:         v.AuthorID,
			Status:           v.Status,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
