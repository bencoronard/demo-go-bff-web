package config

import (
	"time"

	"github.com/bencoronard/demo-go-common-libs/rdb"
	"github.com/caarlos0/env/v11"
	"go.uber.org/fx"
)

type pgDriverCfg struct {
	Host string `mapstructure:"pg.host"`
	Port string `mapstructure:"pg.port"`
	Db   string `mapstructure:"pg.dbname"`
	User string `mapstructure:"pg.user"`
	Pass string `mapstructure:"pg.pass"`
}

type rdbCfg struct {
	MaxOpenConns  int `env:"RDB_CONN_MAX_OPEN"`
	MaxIdleConns  int `env:"RDB_CONN_MAX_IDLE"`
	ConnMaxTTLMs  int `env:"RDB_CONN_MAX_TTL_MSEC"`
	IdleTimeoutMs int `env:"RDB_CONN_IDLE_TIMEOUT_MSEC"`
}

type Properties struct {
	fx.Out
	PgCfg  *rdb.DriverConfig
	RdbCfg *rdb.DBConfig
}

func NewProperties() (Properties, error) {
	pg, err := newPGConfig()
	if err != nil {
		return Properties{}, err
	}

	rdb, err := newRDBConfig()
	if err != nil {
		return Properties{}, err
	}

	return Properties{
		PgCfg:  pg,
		RdbCfg: rdb,
	}, nil
}

func newPGConfig() (*rdb.DriverConfig, error) {
	return &rdb.DriverConfig{}, nil
}

func newRDBConfig() (*rdb.DBConfig, error) {
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
