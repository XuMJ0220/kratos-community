// in internal/registry/registry.go
package registry

import (
	"kratos-community/internal/conf"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
)

var ProviderSet = wire.NewSet(NewRegistry,NewDiscovery)

// NewRegistry 创建一个服务注册器
func NewRegistry(conf *conf.Registry) registry.Registrar {
	// c := conf.Consul
	// if c == nil {
	// 	// 如果配置不存在，可以返回 nil 或者 panic
	// 	// 在我们的设计中，所有服务都需要注册，所以 panic 更合适
	// 	panic("consul config is null")
	// }

	// cli, err := consulAPI.NewClient(&consulAPI.Config{
	// 	Address: c.Address,
	// 	Scheme:  c.Scheme,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// r := consul.New(cli, consul.WithHealthCheck(true))
	// return r
    
    // 逻辑复用
    return newConsulRegistry(conf)
}

// NewDiscovery 创建一个服务发现器
func NewDiscovery(conf *conf.Registry) registry.Discovery {
	return newConsulRegistry(conf)
}

func newConsulRegistry(conf *conf.Registry) *consul.Registry {
	c := conf.Consul
	if c == nil {
		panic("consul config is null")
	}

	cli, err := consulAPI.NewClient(&consulAPI.Config{
		Address: c.Address,
		Scheme:  c.Scheme,
	})
	if err != nil {
		panic(err)
	}

	// consul.WithHealthCheck(true) 是一个好习惯，让 Consul 主动检查服务健康
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}
