package client

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	relationv1 "kratos-community/api/relation/v1"
)



func NewRelationServiceClient(r registry.Discovery) (relationv1.RelationClient, error) {
	endpoint := "discovery:///relation-service"

	// DialInsecure 表示创建一个不使用TLS加密的连接
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint(endpoint), // 设置目标地址，使用服务发现协议
		grpc.WithDiscovery(r),       // 注入我们创建的Consul
	)
	if err != nil {
		return nil, err
	}
	return relationv1.NewRelationClient(conn), nil
}
