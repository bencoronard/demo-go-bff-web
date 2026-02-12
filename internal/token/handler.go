package token

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type TokenHandler struct {
	s tokenService
}

func NewTokenHandler(s tokenService) *TokenHandler {
	return &TokenHandler{s: s}
}

func (h *TokenHandler) GenerateToken(c *echo.Context) error {
	token, err := h.s.issueToken(c.Request().Context())
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, token)
}
