//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"kratos-community/internal/conf"
	"kratos-community/internal/content/biz"
	"kratos-community/internal/content/data"
	"kratos-community/internal/content/service"
	"kratos-community/internal/kafka"
	"kratos-community/internal/registry"
	"kratos-community/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Auth, *conf.Registry, *conf.CacheMode, *conf.Kafka, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, registry.ProviderSet, kafka.ProviderSet, newApp))
}
