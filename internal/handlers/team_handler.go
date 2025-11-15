package handlers

import (
	"log"
	"net/http"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/Grivvus/reviewers/internal/service"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service service.TeamService
}

func NewTeamHandler(ts service.TeamService) *TeamHandler {
	return &TeamHandler{
		service: ts,
	}
}

func (h *TeamHandler) PostTeamAdd(c *gin.Context) {
	var team api.Team

	if err := c.BindJSON(&team); err != nil {
		log.Println(err)
		c.JSON(
			http.StatusBadRequest,
			"Wrong body format was given",
		)
		return
	}

	err := h.service.AddTeam(c.Request.Context(), team)
	if err != nil && err == service.TeamAlreadyExistError {
		log.Println(err)
		c.JSON(
			http.StatusBadRequest,
			newErrorResponse(api.TEAMEXISTS, team.TeamName+" "+err.Error()),
		)
		return
	} else if err != nil {
		log.Println(err)
		c.JSON(
			http.StatusInternalServerError,
			"Unkown server error",
		)
		return
	}

	c.JSON(http.StatusCreated, team)

}

func (h *TeamHandler) GetTeamGet(c *gin.Context, params api.GetTeamGetParams) {
	team, err := h.service.GetTeam(c.Request.Context(), params)

	if err != nil && err == service.ResourceNotFoundError {
		log.Println(err)
		c.JSON(
			http.StatusNotFound,
			newErrorResponse(api.NOTFOUND, err.Error()),
		)
		return
	} else if err != nil {
		log.Println(err)
		c.JSON(
			http.StatusInternalServerError,
			"Unkown server error",
		)
		return
	}

	c.JSON(http.StatusOK, team)
}
