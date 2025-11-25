package entity

type PRReviewers struct {
	ReviewersID   string `db:"reviewer_id" json:"reviewer_id"`
	PullRequestID string `db:"pull_request_id" json:"pull_request_id"`
}
