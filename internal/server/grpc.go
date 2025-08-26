package server

import (
	"kratos-community/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	jwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, auth *conf.Auth, logger log.Logger) *grpc.Server {
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
		"/api.content.v1.Content/CreateArticle",
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
	return srv
}
