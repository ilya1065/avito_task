package entity

import "time"

type PRStatus string

const (
	PROpen   PRStatus = "OPEN"
	PRMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID   string     `db:"pull_request_id" json:"pull_request_id"`
	PullRequestName string     `db:"pull_request_name" json:"pull_request_name"`
	AuthorID        string     `db:"author_id" json:"author_id"`
	Status          PRStatus   `db:"status" json:"status"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	MergedAt        *time.Time `db:"merged_at" json:"merged_at"`
}
