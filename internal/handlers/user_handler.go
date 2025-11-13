package handlers

import (
	"github.com/Grivvus/reviewers/internal/api"
	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetUsersGetReview(c *gin.Context, params api.GetUsersGetReviewParams) {}

func (h *UserHandler) PostUsersSetIsActive(c *gin.Context) {}
