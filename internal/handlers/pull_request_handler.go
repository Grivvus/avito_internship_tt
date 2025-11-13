package handlers

import "github.com/gin-gonic/gin"

type PullRequestHandler struct{}

func NewPullRequestHandler() *PullRequestHandler {
	return &PullRequestHandler{}
}

func (h *PullRequestHandler) PostPullRequestCreate(c *gin.Context) {}

func (h *PullRequestHandler) PostPullRequestMerge(c *gin.Context) {}

func (h *PullRequestHandler) PostPullRequestReassign(c *gin.Context) {}
