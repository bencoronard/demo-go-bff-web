package main

import (
	"github.com/bencoronard/demo-go-bff-web/internal/config"
	"github.com/bencoronard/demo-go-bff-web/internal/permission"
	"github.com/bencoronard/demo-go-bff-web/internal/token"
	"github.com/bencoronard/demo-go-common-libs/http"
	"github.com/bencoronard/demo-go-common-libs/otel"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.NewProperties,
			config.NewDB,
			config.NewJwtIssuer,
			config.NewRouter,
		),
		fx.Provide(
			token.NewTokenService,
			token.NewTokenHandler,
			permission.NewPermissionRepo,
		),
		fx.Provide(
			config.NewResource,
			otel.NewTracerProvider,
			otel.NewMeterProvider,
			otel.NewLoggerProvider,
		),
		fx.Provide(
			http.NewGlobalErrorHandler,
		),
		fx.Invoke(
			config.ConfigureLogger,
			otel.InitOtel,
			http.Router.RegisterMiddlewares,
			http.Router.RegisterRoutes,
		),
		fx.Invoke(http.Start),
	).Run()
}
