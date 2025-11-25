package hendler

import (
	"avito_task/internal/controller/http/v1/hendler/hendlererrors"
	"avito_task/internal/entity"
	"avito_task/internal/service"
	"encoding/json"
	"net/http"
	"time"
)

type PullRequestHendler struct {
	pullRequestrSvc service.PullRequest
}

func NewPullRequestHendler(service service.PullRequest) *PullRequestHendler {
	return &PullRequestHendler{pullRequestrSvc: service}
}

type PullRequestCreate struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type PullRequestCreateResponse struct {
	Pr struct {
		PullRequestId     string          `json:"pull_request_id"`
		PullRequestName   string          `json:"pull_request_name"`
		AuthorId          string          `json:"author_id"`
		Status            entity.PRStatus `json:"status"`
		AssignedReviewers []string        `json:"assigned_reviewers"`
	} `json:"pr"`
}

func (h *PullRequestHendler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var reqDTO PullRequestCreate
	err := json.NewDecoder(r.Body).Decode(&reqDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	PRShort := entity.PullRequestShort{
		PullRequestID:   reqDTO.PullRequestID,
		PullRequestName: reqDTO.PullRequestName,
		AuthorID:        reqDTO.AuthorID,
	}
	PR, assignedReviewers, err := h.pullRequestrSvc.Create(r.Context(), PRShort)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			hendlererrors.WriteError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		}
		if err.Error() == "PR_EXISTS" {
			hendlererrors.WriteError(w, http.StatusConflict, "PR_EXISTS", "PR id already exists")
			return
		}
		hendlererrors.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}
	resp := PullRequestCreateResponse{}
	resp.Pr.PullRequestId = PR.PullRequestID
	resp.Pr.PullRequestName = PR.PullRequestName
	resp.Pr.AuthorId = PR.AuthorID
	resp.Pr.Status = PR.Status
	resp.Pr.AssignedReviewers = assignedReviewers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}

type PullRequestMergeRequest struct {
	PullRequestId string `json:"pull_request_id"`
}

type PullRequestMergeResponse struct {
	Pr struct {
		PullRequestId     string          `json:"pull_request_id"`
		PullRequestName   string          `json:"pull_request_name"`
		AuthorId          string          `json:"author_id"`
		Status            entity.PRStatus `json:"status"`
		AssignedReviewers []string        `json:"assigned_reviewers"`
		MergedAt          *time.Time      `json:"mergedAt"`
	} `json:"pr"`
}

func (h *PullRequestHendler) Merge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var teq PullRequestMergeRequest
	err := json.NewDecoder(r.Body).Decode(&teq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var pr entity.PullRequest

	pr, reviewers, err := h.pullRequestrSvc.Merge(r.Context(), teq.PullRequestId)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			hendlererrors.WriteError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		}
		hendlererrors.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
	var resp PullRequestMergeResponse
	resp.Pr.MergedAt = pr.MergedAt
	resp.Pr.AuthorId = pr.AuthorID
	resp.Pr.AssignedReviewers = reviewers
	resp.Pr.Status = pr.Status
	resp.Pr.PullRequestId = pr.PullRequestID
	resp.Pr.PullRequestName = pr.PullRequestName
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

type ReassignRequest struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

type ReassignResponse struct {
	Pr struct {
		PullRequestId     string          `json:"pull_request_id"`
		PullRequestName   string          `json:"pull_request_name"`
		AuthorId          string          `json:"author_id"`
		Status            entity.PRStatus `json:"status"`
		AssignedReviewers []string        `json:"assigned_reviewers"`
	} `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}

func (h *PullRequestHendler) Reassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var req ReassignRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var pr entity.PullRequest
	var reviewers []string
	var NewReviewer string
	pr, reviewers, NewReviewer, err = h.pullRequestrSvc.Reassign(r.Context(), req.PullRequestId, req.OldReviewerId)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			hendlererrors.WriteError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		}
		if err.Error() == "PR_MERGED" {
			hendlererrors.WriteError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
			return
		}
		if err.Error() == "NO_CANDIDATE" {
			hendlererrors.WriteError(w, http.StatusConflict, "NO_CANDIDATE", "resource not found")
			return
		}
		hendlererrors.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return

	}
	var resp ReassignResponse
	resp.ReplacedBy = NewReviewer
	resp.Pr.AuthorId = pr.AuthorID
	resp.Pr.PullRequestName = pr.PullRequestName
	resp.Pr.AssignedReviewers = reviewers
	resp.Pr.PullRequestId = pr.PullRequestID
	resp.Pr.Status = pr.Status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
