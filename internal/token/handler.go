package token

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

type TokenHandler struct {
	s TokenService
}

func NewTokenHandler(s TokenService) *TokenHandler {
	return &TokenHandler{s: s}
}

func (h *TokenHandler) GenerateToken(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil {
		return err
	}

	sec, err := strconv.ParseUint(c.Param("ttl"), 10, strconv.IntSize)
	if err != nil {
		return err
	}

	token, err := h.s.issueToken(c.Request().Context(), uint(id), time.Duration(sec)*time.Second)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, token)
}
