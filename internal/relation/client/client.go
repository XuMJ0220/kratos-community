package client

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"

	userv1 "kratos-community/api/user/v1"
)

var ProviderSet = wire.NewSet(NewUserServiceClient)

func NewUserServiceClient(r registry.Discovery) (userv1.UserClient, error) {
	endpoint := "discovery:///user-service"

	// DialInsecure 表示创建一个不使用TLS加密的连接
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint(endpoint), // 设置目标地址，使用服务发现协议
		grpc.WithDiscovery(r),       // 注入我们创建的Consul
	)
	if err != nil {
		return nil, err
	}
	return userv1.NewUserClient(conn), nil
}
