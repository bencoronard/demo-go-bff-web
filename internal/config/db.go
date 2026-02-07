package config

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(lc fx.Lifecycle, p *Properties) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.Secret.DB.Host,
		p.Secret.DB.Port,
		p.Secret.DB.User,
		p.Secret.DB.Pass,
		p.Secret.DB.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(p.Env.CP.ConnectionPoolCap)
	sqlDB.SetMaxIdleConns(p.Env.CP.ConnectionPoolIdleMin)
	sqlDB.SetConnMaxLifetime(time.Duration(p.Env.CP.ConnectionTTL) * time.Millisecond)
	sqlDB.SetConnMaxIdleTime(time.Duration(p.Env.CP.ConnectionPoolIdleTimeout) * time.Millisecond)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return sqlDB.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return sqlDB.Close()
		},
	})

	return db, nil
}
