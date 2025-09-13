//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"kratos-community/internal/conf"
	"kratos-community/internal/kafka"
	"kratos-community/internal/notification/biz"
	"kratos-community/internal/notification/client"
	"kratos-community/internal/notification/server"
	"kratos-community/internal/registry"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Registry, *conf.Kafka, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, biz.ProviderSet, registry.ProviderSet, kafka.ProviderSet, client.ProviderSet, newApp))
}
