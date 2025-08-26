package server

import (
	userV1 "kratos-community/api/user/v1"
	"kratos-community/internal/conf"
	"kratos-community/internal/user/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/middleware/selector"
    jwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
    jwtv5 "github.com/golang-jwt/jwt/v5"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, user *service.UserService, auth *conf.Auth, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	// 添加 jwt 中间件
	opts = append(opts, grpc.Middleware(
		selector.Server(
			// 创建 JWT 中间件
			jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
				return []byte(auth.JwtSecret), nil
			}),
		).Path(
		//往这里添加，例如
		//"/api.user.v1.User/RegisterUser"
		).Build(),
	))
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	userV1.RegisterUserServer(srv, user)
	return srv
}
