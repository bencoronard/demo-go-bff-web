package config

import (
	"net/http"

	"github.com/bencoronard/demo-go-bff-web/internal/token"
	xhttp "github.com/bencoronard/demo-go-common-libs/http"
	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/fx"
)

type RouterParams struct {
	fx.In
	p  *Properties
	h  *token.TokenHandler
	eh xhttp.GlobalErrorHandler
}

type router struct {
	p *Properties
	e *echo.Echo
	h *token.TokenHandler
}

func NewRouter(rp RouterParams) xhttp.Router {
	e := echo.New()
	e.HTTPErrorHandler = rp.eh.GetHandler()
	return &router{
		p: rp.p,
		e: e,
		h: rp.h,
	}
}

func (r *router) Port() int {
	return r.p.Env.App.ListenPort
}

func (r *router) Handler() http.Handler {
	return r.e
}

func (r *router) RegisterMiddlewares() {
	r.e.Use(middleware.Recover())
}

func (r *router) RegisterRoutes() {
	api := r.e.Group("/api",
		echootel.NewMiddleware(r.p.Env.App.Name),
		middleware.RequestLogger(),
	)
	api.GET("/token", r.h.GenerateToken)

	act := r.e.Group("/actuator")
	act.GET("/health", healthCheck)
	act.GET("/readiness", readinessCheck)
}

func healthCheck(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "up"})
}

func readinessCheck(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
}
