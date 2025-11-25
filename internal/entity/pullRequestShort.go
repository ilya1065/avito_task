package entity

type PullRequestShort struct {
	PullRequestID   string   `db:"pull_request_id" json:"pull_request_id"`
	PullRequestName string   `db:"pull_request_name" json:"pull_request_name"`
	AuthorID        string   `db:"author_id" json:"author_id"`
	Status          PRStatus `db:"status" json:"status"`
}
