package main

import (
	"github.com/bencoronard/demo-go-bff-web/internal/config"
	"github.com/bencoronard/demo-go-bff-web/internal/token"
	"github.com/bencoronard/demo-go-common-libs/http"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.NewProperties,
			config.NewJwtIssuer,
			config.NewRouter,
		),
		fx.Provide(
			token.NewTokenService,
			token.NewTokenHandler,
		),
		fx.Invoke(
			config.ConfigureLogger,
			http.Router.RegisterMiddlewares,
			http.Router.RegisterRoutes,
		),
		fx.Invoke(http.Start),
	).Run()
}
