package config

import (
	"net"
	"net/http"
	"strconv"

	"github.com/bencoronard/demo-go-bff-web/internal/token"
	"github.com/bencoronard/demo-go-common-libs/actuator"
	"github.com/bencoronard/demo-go-common-libs/server"
	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/fx"
)

type httpServer struct {
	s  *http.Server
	e  *echo.Echo
	a  actuator.Actuator
	th *token.TokenHandler
}

func (h *httpServer) Instance() *http.Server {
	return h.s
}

func (h *httpServer) Configure() error {
	h.e.Use(middleware.Recover())
	h.e.Use(echootel.NewMiddleware(""))
	h.e.Use(middleware.RequestLogger())

	api := h.e.Group("/api/token")
	api.GET("", h.th.GenerateToken)

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
	Prop     server.HttpServerConfig
}

func NewHttpServer(p httpServerParams) server.HttpServer {
	return &httpServer{
		s: &http.Server{
			Addr:              net.JoinHostPort(p.Prop.Host, strconv.Itoa(p.Prop.Port)),
			Handler:           p.Router,
			ReadTimeout:       p.Prop.ReadTimeout,
			ReadHeaderTimeout: p.Prop.ReadHeaderTimeout,
			WriteTimeout:      p.Prop.WriteTimeout,
			IdleTimeout:       p.Prop.IdleTimeout,
			MaxHeaderBytes:    p.Prop.MaxHeaderBytes,
		},
		e:  p.Router,
		a:  p.Actuator,
		th: p.Handler,
	}
}
