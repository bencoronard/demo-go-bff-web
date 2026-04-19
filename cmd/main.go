package main

import (
	"github.com/bencoronard/demo-go-bff-web/internal/config"
	"github.com/bencoronard/demo-go-common-libs/rdb"
	"github.com/bencoronard/demo-go-common-libs/vault"
	"go.uber.org/fx"
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
	).Run()
}
