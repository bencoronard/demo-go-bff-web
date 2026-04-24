package config

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/bencoronard/demo-go-common-libs/actuator"
	"github.com/bencoronard/demo-go-common-libs/jwt"
	"github.com/bencoronard/demo-go-common-libs/rdb"
	"github.com/bencoronard/demo-go-common-libs/server"
	"github.com/bencoronard/demo-go-common-libs/vault"
	"github.com/caarlos0/env/v11"
	"go.uber.org/fx"
)

type rdbCfg struct {
	MaxOpenConn int `env:"RDB_CONN_MAX_OPEN"`
	MaxIdleConn int `env:"RDB_CONN_MAX_IDLE"`
	ConnTTL     int `env:"RDB_CONN_TTL_SEC"`
	IdleTimeout int `env:"RDB_CONN_IDLE_TIMEOUT_SEC"`
}

type pgCfg struct {
	Host string `mapstructure:"pg.host"`
	Port int    `mapstructure:"pg.port"`
	DB   string `mapstructure:"pg.dbname"`
	User string `mapstructure:"pg.user"`
	Pass string `mapstructure:"pg.pass"`
}

type jwtCfg struct {
	Issuer string `env:"SERVICE_NAME"`
	Key    string `mapstructure:"private.key"`
}

type actuatorCfg struct {
	Host                string `env:"SERVER_BIND_HOST"`
	Port                int    `env:"ACTUATOR_PORT"`
	HealthCheckInterval int    `env:"ACTUATOR_HEALTHCHECK_INTERVAL_SEC"`
	HealthCheckTimeout  int    `env:"ACTUATOR_HEALTHCHECK_TIMEOUT_SEC"`
}

type serverCfg struct {
	Host              string `env:"SERVER_BIND_HOST"`
	Port              int    `env:"SERVER_BIND_PORT"`
	ReadTimeout       int    `env:"SERVER_READ_TIMEOUT_SEC"`
	ReadHeaderTimeout int    `env:"SERVER_READ_HEADER_TIMEOUT_SEC"`
	WriteTimeout      int    `env:"SERVER_WRITE_TIMEOUT_SEC"`
	IdleTimeout       int    `env:"SERVER_IDLE_TIMEOUT_SEC"`
	MaxHeaderBytes    int    `env:"SERVER_MAX_HEADER_BYTES"`
}

type properties struct {
	fx.Out
	Rdb rdb.DbConfig
	Pg  rdb.DriverConfig
	Jwt jwt.AsymmIssuerConfig
	Act actuator.Config
	Srv server.HttpServerConfig
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

	act, err := newActuatorCfg()
	if err != nil {
		return properties{}, err
	}

	srv, err := newServerCfg()
	if err != nil {
		return properties{}, err
	}

	return properties{
		Pg:  pg,
		Rdb: rdb,
		Jwt: jwt,
		Act: act,
		Srv: srv,
	}, nil
}

func newPgCfg(vc vault.Client) (rdb.DriverConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c pgCfg
	if err := vc.ReadSecret(ctx, fmt.Sprintf("secret/application/%s", "local"), &c); err != nil {
		return rdb.DriverConfig{}, err
	}

	return rdb.DriverConfig{
		Host:     c.Host,
		Port:     c.Port,
		User:     c.User,
		Password: c.Pass,
		DBName:   c.DB,
		UseSSL:   false,
	}, nil
}

func newRdbCfg() (rdb.DbConfig, error) {
	var c rdbCfg
	if err := env.Parse(&c); err != nil {
		return rdb.DbConfig{}, err
	}
	return rdb.DbConfig{
		MaxOpenConns: c.MaxOpenConn,
		MaxIdleConns: c.MaxIdleConn,
		ConnTTL:      time.Duration(c.ConnTTL) * time.Second,
		IdleTimeout:  time.Duration(c.IdleTimeout) * time.Second,
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

func newActuatorCfg() (actuator.Config, error) {
	var c actuatorCfg
	if err := env.Parse(&c); err != nil {
		return actuator.Config{}, err
	}
	return actuator.Config{
		HealthCheckInterval: time.Duration(c.HealthCheckInterval) * time.Second,
		HealthCheckTimeout:  time.Duration(c.HealthCheckTimeout) * time.Second,
	}, nil
}

func newServerCfg() (server.HttpServerConfig, error) {
	var c serverCfg
	if err := env.Parse(&c); err != nil {
		return server.HttpServerConfig{}, err
	}
	return server.HttpServerConfig{
		Host:              c.Host,
		Port:              c.Port,
		ReadTimeout:       time.Duration(c.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(c.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(c.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(c.IdleTimeout) * time.Second,
		MaxHeaderBytes:    c.MaxHeaderBytes,
	}, nil
}
