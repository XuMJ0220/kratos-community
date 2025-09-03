package client

import (
	"context"
	contentv1 "kratos-community/api/content/v1"
	userv1 "kratos-community/api/user/v1"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

// ProviderSet 是客户端的依赖注入提供者集合
var ProviderSet = wire.NewSet(NewUserServiceClient, NewContentServiceClient)

func NewUserServiceClient(r registry.Discovery) (userv1.UserClient, error) {
	// discovery:///user-service 是kratos定义的服务发现协议
	// "user-service" 是我们在user-service 的main.go中为App定义的Name
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

func NewContentServiceClient(r registry.Discovery) (contentv1.ContentClient, error) {
	endpoint := "discovery:///content-service"

	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		return nil, err
	}
	return contentv1.NewContentClient(conn), nil
}
