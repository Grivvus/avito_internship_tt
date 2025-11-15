package handlers

import (
	"net/http"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/Grivvus/reviewers/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{
		service: us,
	}
}

func (h *UserHandler) PostUsersSetIsActive(c *gin.Context) {
	var model api.PostUsersSetIsActiveJSONBody

	if err := c.BindJSON(&model); err != nil {
		c.JSON(
			http.StatusBadRequest,
			"Wrong body format was given",
		)
		return
	}

	responseModel, err := h.service.SetIsActive(c.Request.Context(), model)
	if err != nil && err == service.ResourceNotFoundError {
		c.JSON(
			http.StatusNotFound,
			newErrorResponse(api.NOTFOUND, err.Error()),
		)
		return
	} else if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			"Unkown server error",
		)
		return
	}

	c.JSON(
		http.StatusOK,
		responseModel,
	)
}

func (h *UserHandler) GetUsersGetReview(
	c *gin.Context,
	params api.GetUsersGetReviewParams,
) {
	response, err := h.service.UserReviews(c.Request.Context(), params)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			"Unkown server error",
		)
		return
	}

	c.JSON(http.StatusOK, response)
}
