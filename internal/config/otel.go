package config

import (
	"go.opentelemetry.io/otel/sdk/resource"
)

func NewResource() (*resource.Resource, error) {
	return resource.Environment(), nil
}
