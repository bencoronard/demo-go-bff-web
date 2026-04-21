package config

import (
	"net/http"

	"github.com/bencoronard/demo-go-bff-web/internal/token"
	"github.com/bencoronard/demo-go-common-libs/actuator"
	"github.com/bencoronard/demo-go-common-libs/server"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/fx"
)

type httpServer struct {
	e *echo.Echo
	a actuator.Actuator
	h *token.TokenHandler
}

func (h *httpServer) Instance() *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: h.e,
	}
}

func (h *httpServer) Configure() error {
	h.e.Use(middleware.Recover())

	api := h.e.Group("/api/tokens")
	api.GET("", h.h.GenerateToken)

	act := h.e.Group("/actuator")
	act.GET("/liveness", func(c *echo.Context) error {
		if live := h.a.Liveness(); !live {
			return c.NoContent(http.StatusServiceUnavailable)
		}
		return c.NoContent(http.StatusNoContent)
	})
	act.GET("/readiness", func(c *echo.Context) error {
		if ready := h.a.Readiness(); !ready {
			return c.NoContent(http.StatusServiceUnavailable)
		}
		return c.NoContent(http.StatusNoContent)
	})

	return nil
}

type httpServerParams struct {
	fx.In
	Router   *echo.Echo
	Actuator actuator.Actuator
	Handler  *token.TokenHandler
}

func NewHttpServer(p httpServerParams) server.HttpServer {
	return &httpServer{
		e: p.Router,
		a: p.Actuator,
		h: p.Handler,
	}
}
