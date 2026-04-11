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
	MaxOpenConns  int `env:"RDB_CONN_MAX_OPEN"`
	MaxIdleConns  int `env:"RDB_CONN_MAX_IDLE"`
	ConnMaxTTLMs  int `env:"RDB_CONN_MAX_TTL_MSEC"`
	IdleTimeoutMs int `env:"RDB_CONN_IDLE_TIMEOUT_MSEC"`
}

type pgCfg struct {
	Host string `mapstructure:"pg.host"`
	Port string `mapstructure:"pg.port"`
	Db   string `mapstructure:"pg.dbname"`
	User string `mapstructure:"pg.user"`
	Pass string `mapstructure:"pg.pass"`
}

type Properties struct {
	fx.Out
	RdbCfg *rdb.DBConfig
	PgCfg  *rdb.DriverConfig
}

type PropParams struct {
	fx.In
	Vc vault.Client
}

func NewProperties(p PropParams) (Properties, error) {
	rdb, err := newRdbCfg()
	if err != nil {
		return Properties{}, err
	}

	pg, err := newPgCfg(p.Vc)
	if err != nil {
		return Properties{}, err
	}

	return Properties{
		PgCfg:  pg,
		RdbCfg: rdb,
	}, nil
}

func newPgCfg(vc vault.Client) (*rdb.DriverConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c pgCfg
	if err := vc.ReadSecret(ctx, fmt.Sprintf("secret/application/%s", "dev"), &c); err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return nil, err
	}

	return &rdb.DriverConfig{
		Host:     c.Host,
		Port:     port,
		User:     c.User,
		Password: c.Pass,
		DBName:   c.Db,
		UseSSL:   false,
	}, nil
}

func newRdbCfg() (*rdb.DBConfig, error) {
	var cfg rdbCfg
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &rdb.DBConfig{
		MaxOpenConns: cfg.MaxOpenConns,
		MaxIdleConns: cfg.MaxIdleConns,
		ConnMaxTTL:   time.Duration(cfg.ConnMaxTTLMs) * time.Millisecond,
		IdleTimeout:  time.Duration(cfg.IdleTimeoutMs) * time.Millisecond,
	}, nil
}
