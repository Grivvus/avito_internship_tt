package handlers

import (
	"github.com/Grivvus/reviewers/internal/api"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct{}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{}
}

func (h *TeamHandler) PostTeamAdd(c *gin.Context) {}

func (h *TeamHandler) GetTeamGet(c *gin.Context, params api.GetTeamGetParams) {}
