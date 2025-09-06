//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"kratos-community/internal/conf"
	"kratos-community/internal/registry"
	"kratos-community/internal/relation/biz"
	"kratos-community/internal/relation/client"
	"kratos-community/internal/relation/data"
	"kratos-community/internal/relation/service"
	"kratos-community/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, *conf.Auth, *conf.Registry, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, client.ProviderSet, registry.ProviderSet, newApp))
}
