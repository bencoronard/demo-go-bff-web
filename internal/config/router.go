package config

import (
	"net/http"

	"github.com/bencoronard/demo-go-bff-web/internal/token"
	xhttp "github.com/bencoronard/demo-go-common-libs/http"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type router struct {
	port int
	e    *echo.Echo
	h    *token.TokenHandler
}

func NewRouter(p *Properties, h *token.TokenHandler) xhttp.Router {
	e := echo.New()
	e.HTTPErrorHandler = xhttp.GlobalErrorHandler(nil)
	return &router{
		port: p.Env.App.ListenPort,
		e:    e,
		h:    h,
	}
}

func (r *router) Port() int {
	return r.port
}

func (r *router) Handler() http.Handler {
	return r.e
}

func (r *router) RegisterMiddlewares() {
	r.e.Use(middleware.Recover())
}

func (r *router) RegisterRoutes() {
	api := r.e.Group("/api", middleware.RequestLogger())
	api.GET("/token", r.h.GenerateToken)

	act := r.e.Group("/actuator")
	act.GET("/health", func(c *echo.Context) error { return c.JSON(http.StatusOK, map[string]string{"status": "up"}) })
	act.GET("/readiness", func(c *echo.Context) error { return c.JSON(http.StatusOK, map[string]string{"status": "ready"}) })
}
