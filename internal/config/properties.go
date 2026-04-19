package config

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bencoronard/demo-go-common-libs/jwt"
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

type jwtCfg struct {
	Issuer string `env:"APP_NAME"`
	Key    string `mapstructure:"private.key"`
}

type properties struct {
	fx.Out
	Rdb rdb.DbConfig
	Pg  rdb.DriverConfig
	Jwt jwt.AsymmIssuerConfig
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

	jwt, err := newJwtIssuerCfg(p.Vc)
	if err != nil {
		return properties{}, err
	}

	return properties{
		Pg:  pg,
		Rdb: rdb,
		Jwt: jwt,
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

func newJwtIssuerCfg(vc vault.Client) (jwt.AsymmIssuerConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c jwtCfg
	if err := env.Parse(&c); err != nil {
		return jwt.AsymmIssuerConfig{}, err
	}
	if err := vc.ReadSecret(ctx, "secret/bff-web", &c); err != nil {
		return jwt.AsymmIssuerConfig{}, err
	}

	block, _ := pem.Decode([]byte(c.Key))
	if block == nil {
		return jwt.AsymmIssuerConfig{}, errors.New("failed to parse private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return jwt.AsymmIssuerConfig{}, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return jwt.AsymmIssuerConfig{}, errors.New("not an RSA private key")
	}

	return jwt.AsymmIssuerConfig{
		Issuer: c.Issuer,
		Key:    rsaKey,
	}, nil
}
