package service

import (
	"context"

	contentv1 "kratos-community/api/content/v1"
	pb "kratos-community/api/gateway/v1"
	userv1 "kratos-community/api/user/v1"

	"github.com/go-kratos/kratos/v2/errors"
	jwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// 自定义错误类型
var (
	ErrInternalServer = errors.InternalServer("Err_INTERNAL_SERVER", "服务器出错")
	ErrAuthorized     = errors.Unauthorized("UNAUTHORIZED", "user_id missing in claims")
)

type GatewayService struct {
	pb.UnimplementedGatewayServer
	userClient    userv1.UserClient
	contentClient contentv1.ContentClient
}

func NewGatewayService(userClient userv1.UserClient, contentClient contentv1.ContentClient) *GatewayService {
	return &GatewayService{userClient: userClient, contentClient: contentClient}
}

func (s *GatewayService) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginReply, error) {
	//return &userv1.LoginReply{}, nil
	return s.userClient.Login(ctx, req)
}
func (s *GatewayService) RegisterUser(ctx context.Context, req *userv1.RegisterUserRequest) (*userv1.RegisterUserReply, error) {
	// return &userv1.RegisterUserReply{}, nil
	return s.userClient.RegisterUser(ctx, req)
}
func (s *GatewayService) CreateArticle(ctx context.Context, req *contentv1.CreateArticleRequest) (*contentv1.CreateArticleReply, error) {
	//return &contentv1.CreateArticleReply{}, nil

	// // 从 jwt 中获取user_id
	// claims, ok := jwt.FromContext(ctx)
	// if !ok {
	// 	return nil, ErrInternalServer
	// }
	// // 将cliams断言为我们熟悉的map类型
	// mapClaims, ok := claims.(jwtv5.MapClaims)
	// if !ok {
	// 	return nil, ErrInternalServer
	// }
	// // 取出user_id
	// userId, ok := mapClaims["user_id"].(float64)
	// if !ok {
	// 	return nil, ErrAuthorized
	// }
	// authorId := uint64(userId)
	authorId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	req.AuthorId = authorId

	return s.contentClient.CreateArticle(ctx, req)
}

func (s *GatewayService) GetArticle(ctx context.Context, req *contentv1.GetArticleRequest) (*contentv1.GetArticleReply, error) {
	return s.contentClient.GetArticle(ctx, req)
}

func (s *GatewayService) UpdateArticle(ctx context.Context, req *contentv1.UpdateArticleRequest) (*contentv1.UpdateArticleReply, error) {
	authorId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	req.AuthorId = authorId
	return s.contentClient.UpdateArticle(ctx, req)
}

func (s *GatewayService) DeleteArticle(ctx context.Context, req *contentv1.DeleteArticleRequest) (*contentv1.DeleteArticlReply, error) {
	authorId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	req.AuthorId = authorId
	return s.contentClient.DeleteArticle(ctx, req)
}

// getUserId 获取从Token中携带的id
func getUserId(ctx context.Context) (uint64, error) {
	// 从 jwt 中获取user_id
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return 0, ErrInternalServer
	}
	// 将cliams断言为我们熟悉的map类型
	mapClaims, ok := claims.(jwtv5.MapClaims)
	if !ok {
		return 0, ErrInternalServer
	}
	// 取出user_id
	userId, ok := mapClaims["user_id"].(float64)
	if !ok {
		return 0, ErrAuthorized
	}
	id := uint64(userId)
	return id, nil
}
