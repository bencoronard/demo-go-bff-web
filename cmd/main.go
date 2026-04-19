package main

import (
	"log/slog"

	"github.com/bencoronard/demo-go-bff-web/internal/config"
	"github.com/bencoronard/demo-go-common-libs/jwt"
	"github.com/bencoronard/demo-go-common-libs/rdb"
	"github.com/bencoronard/demo-go-common-libs/vault"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		fx.Provide(
			vault.NewClient,
		),
		fx.Provide(
			config.NewProperties,
		),
		fx.Provide(
			rdb.NewPgDriver,
			rdb.NewDb,
		),
		fx.Provide(
			jwt.NewAsymmIssuer,
		),
		fx.Invoke(func(db *gorm.DB, jwt jwt.Issuer) {
			slog.Info("Application started")
		}),
	).Run()
}
