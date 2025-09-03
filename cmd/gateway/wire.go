//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"kratos-community/internal/conf"
	"kratos-community/internal/gateway/client"
	"kratos-community/internal/gateway/service"
	"kratos-community/internal/registry"
	"kratos-community/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Registry, *conf.Auth,log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet,service.ProviderSet, client.ProviderSet, registry.ProviderSet, newApp))
}
