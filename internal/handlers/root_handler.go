package handlers

import "github.com/Grivvus/reviewers/internal/api"

type RootHandler struct {
	*UserHandler
	*TeamHandler
	*PullRequestHandler
}

func NewRootHandler(
	uh *UserHandler,
	th *TeamHandler,
	prh *PullRequestHandler,
) *RootHandler {
	return &RootHandler{
		UserHandler:        uh,
		TeamHandler:        th,
		PullRequestHandler: prh,
	}
}

func newErrorResponse(code api.ErrorResponseErrorCode, message string) api.ErrorResponse {
	return api.ErrorResponse{
		Error: struct {
			Code    api.ErrorResponseErrorCode `json:"code"`
			Message string                     `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}
}
