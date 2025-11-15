package models

import "github.com/Grivvus/reviewers/internal/api"

type UserReviewsResponse struct {
	UserId       string                 `json:"user_id"`
	PullRequests []api.PullRequestShort `json:"pull_requests"`
}
