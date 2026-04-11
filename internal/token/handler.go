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
	dto := issueTokenDTO{
		Id:  c.Param("id"),
		Ttl: c.Param("ttl"),
	}

	if err := c.Validate(&dto); err != nil {
		return err
	}

	id, err := strconv.ParseUint(dto.Id, 10, strconv.IntSize)
	if err != nil {
		return err
	}

	ttl, err := strconv.ParseUint(dto.Ttl, 10, strconv.IntSize)
	if err != nil {
		return err
	}

	token, err := h.s.IssueToken(c.Request().Context(), uint(id), time.Duration(ttl)*time.Second)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, token)
}
