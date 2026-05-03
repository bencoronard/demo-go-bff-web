package config

import (
	"net"
	"net/http"
	"strconv"

	"github.com/bencoronard/demo-go-bff-web/internal/token"
	"github.com/bencoronard/demo-go-common-libs/server"
	"github.com/labstack/echo/v5"
	"go.uber.org/fx"
)

type httpServer struct {
	s  *http.Server
	e  *echo.Echo
	th *token.TokenHandler
}

func (h *httpServer) Instance() *http.Server {
	return h.s
}

func (h *httpServer) Configure() error {
	api := h.e.Group("/token")
	api.GET("", h.th.GenerateToken)
	return nil
}

type httpServerParams struct {
	fx.In
	Router  *echo.Echo
	Handler *token.TokenHandler
	Prop    server.HTTPServerConfig
}

func NewHttpServer(p httpServerParams) server.HTTPServer {
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
		th: p.Handler,
	}
}
