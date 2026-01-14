package handler

import (
	"errors"
	"net/http"

	"upwork-test/internal/application/usecase"
	"upwork-test/internal/delivery/http/request"
	"upwork-test/internal/delivery/http/response"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	getCategoryOverviewUseCase *usecase.GetCategoryOverview
}

func NewCategoryHandler(
	getCategoryOverviewUseCase *usecase.GetCategoryOverview,
) *CategoryHandler {
	return &CategoryHandler{
		getCategoryOverviewUseCase: getCategoryOverviewUseCase,
	}
}

func (h *CategoryHandler) GetOverview(c *gin.Context) {
	traceID, _ := c.Get("trace_id")

	var req request.GetCategoryOverviewRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid category",
			traceID.(string),
		))
		return
	}

	result, err := h.getCategoryOverviewUseCase.Execute(c.Request.Context(), req.Category)
	if err != nil {
		if errors.Is(err, usecase.ErrCategoryOverviewNotFound) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(
				http.StatusNotFound,
				"Category overview not found",
				traceID.(string),
			))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to get category overview",
			traceID.(string),
		))
		return
	}

	apiResponse := response.FromCategoryOverviewDTO(result)

	c.JSON(http.StatusOK, apiResponse)
}
