package config

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/bencoronard/demo-go-bff-web/internal/token"
	xhttp "github.com/bencoronard/demo-go-common-libs/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type router struct {
	port int
	e    *echo.Echo
	h    *token.TokenHandler
}

func NewRouter(h *token.TokenHandler, p *Properties) xhttp.Router {
	return &router{
		port: p.Env.App.ListenPort,
		e:    echo.New(),
		h:    h,
	}
}

func (r *router) ListeningPort() int {
	return r.port
}

func (r *router) Listen(port int) (net.Listener, error) {
	addr := fmt.Sprintf(":%d", port)
	return net.Listen("tcp", addr)
}

func (r *router) Serve(l net.Listener) error {
	return r.e.Server.Serve(l)
}

func (r *router) Shutdown(ctx context.Context) error {
	return r.e.Shutdown(ctx)
}

func (r *router) RegisterMiddlewares() {
	r.e.Use(middleware.Recover())
}

func (r *router) RegisterRoutes() {
	api := r.e.Group("/api", middleware.RequestLogger())
	api.GET("/token", r.h.GenerateToken)

	act := r.e.Group("/actuator")
	act.GET("/health", func(c echo.Context) error { return c.JSON(http.StatusOK, map[string]string{"status": "up"}) })
	act.GET("/readiness", func(c echo.Context) error { return c.JSON(http.StatusOK, map[string]string{"status": "ready"}) })
}
