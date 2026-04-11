package main

import (
	"github.com/bencoronard/demo-go-bff-web/internal/config"
	"github.com/bencoronard/demo-go-common-libs/rdb"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.NewProperties,
		),
		fx.Provide(
			rdb.NewPGDriver,
			rdb.NewDB,
		),
	).Run()
}
