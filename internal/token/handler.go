package token

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type TokenHandler struct {
	s tokenService
}

func NewTokenHandler(s tokenService) *TokenHandler {
	return &TokenHandler{s: s}
}

func (h *TokenHandler) GenerateToken(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil {
		return err
	}

	token, err := h.s.issueToken(c.Request().Context(), uint(id))
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, token)
}
