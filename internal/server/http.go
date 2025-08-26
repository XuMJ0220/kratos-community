package server

import (
	"kratos-community/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	jwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, auth *conf.Auth, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	}

	// 添加 jwt 中间件
	opts = append(opts, http.Middleware(
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
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	return srv
}
