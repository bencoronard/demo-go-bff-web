package config

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
)

func ConfigureLogger(lp *log.LoggerProvider) {
	handler := otelslog.NewHandler(
		"",
		otelslog.WithLoggerProvider(lp),
	)

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
