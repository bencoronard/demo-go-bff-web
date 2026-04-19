package config

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bencoronard/demo-go-common-libs/rdb"
	"github.com/bencoronard/demo-go-common-libs/vault"
	"github.com/caarlos0/env/v11"
	"go.uber.org/fx"
)

type rdbCfg struct {
	MaxOpenConn int `env:"RDB_CONN_MAX_OPEN"`
	MaxIdleConn int `env:"RDB_CONN_MAX_IDLE"`
	ConnTTL     int `env:"RDB_CONN_TTL_MILLISEC"`
	IdleTimeout int `env:"RDB_CONN_IDLE_TIMEOUT_MILLISEC"`
}

type pgCfg struct {
	Host string `mapstructure:"pg.host"`
	Port string `mapstructure:"pg.port"`
	DB   string `mapstructure:"pg.dbname"`
	User string `mapstructure:"pg.user"`
	Pass string `mapstructure:"pg.pass"`
}

type properties struct {
	fx.Out
	Rdb rdb.DbConfig
	Pg  rdb.DriverConfig
}

type propParams struct {
	fx.In
	Vc vault.Client
}

func NewProperties(p propParams) (properties, error) {
	rdb, err := newRdbCfg()
	if err != nil {
		return properties{}, err
	}

	pg, err := newPgCfg(p.Vc)
	if err != nil {
		return properties{}, err
	}

	return properties{
		Pg:  pg,
		Rdb: rdb,
	}, nil
}

func newPgCfg(vc vault.Client) (rdb.DriverConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c pgCfg
	if err := vc.ReadSecret(ctx, fmt.Sprintf("secret/application/%s", "local"), &c); err != nil {
		return rdb.DriverConfig{}, err
	}

	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return rdb.DriverConfig{}, err
	}

	return rdb.DriverConfig{
		Host:     c.Host,
		Port:     port,
		User:     c.User,
		Password: c.Pass,
		DBName:   c.DB,
		UseSSL:   false,
	}, nil
}

func newRdbCfg() (rdb.DbConfig, error) {
	var cfg rdbCfg
	if err := env.Parse(&cfg); err != nil {
		return rdb.DbConfig{}, err
	}
	return rdb.DbConfig{
		MaxOpenConns: cfg.MaxOpenConn,
		MaxIdleConns: cfg.MaxIdleConn,
		ConnTTL:      time.Duration(cfg.ConnTTL) * time.Millisecond,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Millisecond,
	}, nil
}
