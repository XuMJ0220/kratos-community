package service

import (
	"context"

	pb "kratos-community/api/content/v1"
	"kratos-community/internal/content/biz"

	"github.com/go-kratos/kratos/v2/errors"
	jwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type ContentService struct {
	pb.UnimplementedContentServer

	uc *biz.ContentUsecase
}

func NewContentService(uc *biz.ContentUsecase) *ContentService {
	return &ContentService{uc: uc}
}

func (s *ContentService) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleReply, error) {

	// 从 jwt 中获取user_id
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized("UNAUTHORIZED", "token missing or invalid")
	}
	// 将cliams断言为我们熟悉的map类型
	mapClaims, ok := claims.(jwtv5.MapClaims)
	if !ok {
		return nil, errors.Unauthorized("UNAUTHORIZED", "invalid claims type")
	}
	// 取出user_id
	userId, ok := mapClaims["user_id"].(float64)
	if !ok {
		return nil, errors.Unauthorized("UNAUTHORIZED", "user_id missing in claims")
	}
	authorId := uint64(userId)
	
	article, err := s.uc.CreArticle(ctx, authorId, req.Title, req.Content)
	if err != nil {
		return nil, err
	}
	return &pb.CreateArticleReply{
		Article: &pb.Article{
			Id:        article.Id,
			Title:     article.Title,
			Content:   article.Content,
			AuthorId:  article.AuthorId,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		},
	}, nil
}
