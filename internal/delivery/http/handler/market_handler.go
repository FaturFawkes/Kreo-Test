package handler

import (
	"errors"
	"net/http"

	"upwork-test/internal/application/usecase"
	"upwork-test/internal/delivery/http/request"
	"upwork-test/internal/delivery/http/response"

	"github.com/gin-gonic/gin"
)

type MarketHandler struct {
	listMarketsUseCase      *usecase.ListMarkets
	getMarketDetailsUseCase *usecase.GetMarketDetails
}

func NewMarketHandler(
	listMarketsUseCase *usecase.ListMarkets,
	getMarketDetailsUseCase *usecase.GetMarketDetails,
) *MarketHandler {
	return &MarketHandler{
		listMarketsUseCase:      listMarketsUseCase,
		getMarketDetailsUseCase: getMarketDetailsUseCase,
	}
}

func (h *MarketHandler) ListMarkets(c *gin.Context) {
	traceID, _ := c.Get("trace_id")

	var req request.ListMarketsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid category",
			traceID.(string),
		))
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid query parameters",
			traceID.(string),
		))
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	result, err := h.listMarketsUseCase.Execute(c.Request.Context(), req.Category, req.Page, req.Limit, req.Status)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidPage) || errors.Is(err, usecase.ErrInvalidLimit) {
			c.JSON(http.StatusBadRequest, response.NewErrorResponse(
				http.StatusBadRequest,
				err.Error(),
				traceID.(string),
			))
			return
		}

		if errors.Is(err, usecase.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(
				http.StatusNotFound,
				"Category not found",
				traceID.(string),
			))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			http.StatusInternalServerError,
			"Internal server error",
			traceID.(string),
		))
		return
	}

	c.Header("Cache-Control", "public, max-age=300")
	c.JSON(http.StatusOK, response.FromMarketListDTO(result))
}

func (h *MarketHandler) GetMarketDetails(c *gin.Context) {
	traceID, _ := c.Get("trace_id")

	ticker := c.Param("ticker")
	if ticker == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			http.StatusBadRequest,
			"Ticker is required",
			traceID.(string),
		))
		return
	}

	result, err := h.getMarketDetailsUseCase.Execute(c.Request.Context(), ticker)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidTicker) {
			c.JSON(http.StatusBadRequest, response.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid ticker format",
				traceID.(string),
			))
			return
		}

		if errors.Is(err, usecase.ErrMarketNotFound) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(
				http.StatusNotFound,
				"Market not found",
				traceID.(string),
			))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			http.StatusInternalServerError,
			"Internal server error",
			traceID.(string),
		))
		return
	}

	c.Header("Cache-Control", "public, max-age=30")

	statusCode := http.StatusOK
	if result.IsPartial {
		statusCode = http.StatusPartialContent
	}

	c.JSON(statusCode, response.FromMarketDetailDTO(result))
}
