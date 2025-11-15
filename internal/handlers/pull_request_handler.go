package handlers

import (
	"log"
	"net/http"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/Grivvus/reviewers/internal/service"
	"github.com/gin-gonic/gin"
)

type PullRequestHandler struct {
	service service.PullReqeustService
}

func NewPullRequestHandler(prService service.PullReqeustService) *PullRequestHandler {
	return &PullRequestHandler{
		service: prService,
	}
}

func (h *PullRequestHandler) PostPullRequestCreate(c *gin.Context) {
	var prCreate api.PostPullRequestCreateJSONBody

	if err := c.BindJSON(&prCreate); err != nil {
		log.Println(err)
		c.JSON(
			http.StatusBadRequest,
			"Wrong body format was given",
		)
		return
	}

	response, err := h.service.Create(c.Request.Context(), prCreate)
	if err != nil {
		log.Println(err)
		switch err {
		case service.ResourceNotFoundError:
			c.JSON(
				http.StatusNotFound,
				newErrorResponse(api.NOTFOUND, err.Error()),
			)
		case service.PRAlreadyExistError:
			c.JSON(
				http.StatusConflict,
				newErrorResponse(api.PREXISTS, "PR "+prCreate.AuthorId+" "+err.Error()),
			)
		default:
			c.JSON(
				http.StatusConflict,
				"Unkown server error",
			)
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *PullRequestHandler) PostPullRequestMerge(c *gin.Context) {
	var prToMerge api.PostPullRequestMergeJSONBody

	if err := c.BindJSON(&prToMerge); err != nil {
		log.Println(err)
		c.JSON(
			http.StatusBadRequest,
			"Wrong body format was given",
		)
		return
	}

	response, err := h.service.Merge(c.Request.Context(), prToMerge)
	if err != nil {
		log.Println(err)
		if err == service.ResourceNotFoundError {
			c.JSON(
				http.StatusNotFound,
				newErrorResponse(api.NOTFOUND, err.Error()),
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				"Unkown server error",
			)
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PullRequestHandler) PostPullRequestReassign(c *gin.Context) {
	var prToReassign api.PostPullRequestReassignJSONBody

	if err := c.BindJSON(&prToReassign); err != nil {
		log.Println(err)
		c.JSON(
			http.StatusBadRequest,
			"Wrong body format was given",
		)
		return
	}

	response, err := h.service.Reassign(c.Request.Context(), prToReassign)
	if err != nil {
		log.Println(err)
		switch err {
		case service.ResourceNotFoundError:
			c.JSON(
				http.StatusNotFound,
				newErrorResponse(api.NOTFOUND, err.Error()),
			)
		case service.CantReassignOnMergedPRError, service.NoCandidatesError, service.ReviewerNotAssignedError:
			c.JSON(
				http.StatusConflict,
				err.Error(),
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				"Unkown server error",
			)
		}
		return
	}
	c.JSON(http.StatusOK, response)
}
