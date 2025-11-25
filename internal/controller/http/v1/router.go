package v1

import (
	"avito_task/internal/controller/http/v1/hendler"
	"net/http"
)

func NewRouter(hendler *hendler.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/team/add", hendler.TeamHandler.CreateTeam)
	mux.HandleFunc("/team/get", hendler.TeamHandler.GetTeam)
	mux.HandleFunc("/users/setIsActive", hendler.UserHandler.SetIsActive)
	mux.HandleFunc("/users/getReview", hendler.UserHandler.GetReview)
	mux.HandleFunc("/pullRequest/create", hendler.PrHandler.Create)
	mux.HandleFunc("/pullRequest/merge", hendler.PrHandler.Merge)
	mux.HandleFunc("/pullRequest/reassign", hendler.PrHandler.Reassign)
	return mux
}
