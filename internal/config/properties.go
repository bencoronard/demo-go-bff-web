package config

import (
	"context"
	"fmt"
	"time"

	"github.com/bencoronard/demo-go-common-libs/vault"
	"github.com/caarlos0/env/v11"
	"go.uber.org/fx"
)

type Properties struct {
	Env    envCfg
	Secret secretCfg
}

type envCfg struct {
	App   appCfg
	Vault vaultCfg
	OTEL  otelCfg
	CP    cpCfg
}

type secretCfg struct {
	DB dbCfg `mapstructure:",squash"`
}

func NewProperties(lc fx.Lifecycle) (*Properties, error) {
	var e envCfg
	if err := env.Parse(&e); err != nil {
		return nil, err
	}

	vc, err := vault.NewTokenClient(lc, e.Vault.URI, e.Vault.Token)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var s secretCfg
	if err := vc.ReadSecret(ctx, fmt.Sprintf("secret/application/%s", e.App.Environment), &s); err != nil {
		return nil, err
	}

	return &Properties{
		Env:    e,
		Secret: s,
	}, nil
}

type appCfg struct {
	ListenPort    int    `env:"APP_LISTEN_PORT"`
	Environment   string `env:"APP_ENVIRONMENT"`
	PublicKeyPath string `env:"APP_PUBLIC_KEY_PATH"`
}

type vaultCfg struct {
	URI   string `env:"VAULT_URI"`
	Token string `env:"VAULT_TOKEN"`
}

type otelCfg struct {
	MetricsEndpoint           string  `env:"OTEL_COL_METRICS_ENDPOINT"`
	TracesEndpoint            string  `env:"OTEL_COL_TRACES_ENDPOINT"`
	LogsEndpoint              string  `env:"OTEL_COL_LOGS_ENDPOINT"`
	MetricsSamplingFreqInMin  string  `env:"OTEL_METRICS_SAMPLING_FREQ_IN_MIN"`
	TracesSamplingProbability float64 `env:"OTEL_TRACES_SAMPLING_PROBABILITY"`
}

type cpCfg struct {
	ConnectionPoolCap         int `env:"CRUD_API_DB_CP_CAP"`
	ConnectionPoolIdleMin     int `env:"CRUD_API_DB_CP_IDLE_MIN"`
	ConnectionPoolIdleTimeout int `env:"CRUD_API_DB_CP_IDLE_TIMEOUT"`
	ConnectionTimeout         int `env:"CRUD_API_DB_CP_CONN_TIMEOUT"`
	ConnectionTTL             int `env:"CRUD_API_DB_CP_CONN_TTL"`
}

type dbCfg struct {
	Host string `mapstructure:"pg.host"`
	Port string `mapstructure:"pg.port"`
	Name string `mapstructure:"pg.dbname"`
	User string `mapstructure:"pg.user"`
	Pass string `mapstructure:"pg.pass"`
}
