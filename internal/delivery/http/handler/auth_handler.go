package handler

import (
	"errors"
	"net/http"
	"upwork-test/internal/application/dto"
	"upwork-test/internal/application/usecase"
	"upwork-test/internal/delivery/http/request"
	"upwork-test/internal/delivery/http/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authenticateUseCase *usecase.Authenticate
}

func NewAuthHandler(authenticateUseCase *usecase.Authenticate) *AuthHandler {
	return &AuthHandler{
		authenticateUseCase: authenticateUseCase,
	}
}

func (h *AuthHandler) IssueToken(c *gin.Context) {
	traceID, _ := c.Get("trace_id")

	var req request.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			http.StatusBadRequest,
			"Invalid request body",
			traceID.(string),
		))
		return
	}

	tokenResponse, err := h.authenticateUseCase.Execute(&dto.AuthRequest{
		APIKey: req.APIKey,
	})

	if err != nil {
		if errors.Is(err, usecase.ErrInvalidAPIKey) {
			c.JSON(http.StatusUnauthorized, response.NewErrorResponse(
				http.StatusUnauthorized,
				"Invalid API key",
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

	c.JSON(http.StatusOK, response.NewTokenResponse(
		tokenResponse.Token,
		tokenResponse.ExpiresAt,
	))
}
