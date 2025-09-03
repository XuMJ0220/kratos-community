package service

import (
	"context"

	pb "kratos-community/api/content/v1"
	"kratos-community/internal/content/biz"
)

type ContentService struct {
	pb.UnimplementedContentServer

	uc *biz.ContentUsecase
}

func NewContentService(uc *biz.ContentUsecase) *ContentService {
	return &ContentService{uc: uc}
}

func (s *ContentService) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleReply, error) {

	// // 从 jwt 中获取user_id
	// claims, ok := jwt.FromContext(ctx)
	// if !ok {
	// 	return nil, errors.Unauthorized("UNAUTHORIZED", "token missing or invalid")
	// }
	// // 将cliams断言为我们熟悉的map类型
	// mapClaims, ok := claims.(jwtv5.MapClaims)
	// if !ok {
	// 	return nil, errors.Unauthorized("UNAUTHORIZED", "invalid claims type")
	// }
	// // 取出user_id
	// userId, ok := mapClaims["user_id"].(float64)
	// if !ok {
	// 	return nil, errors.Unauthorized("UNAUTHORIZED", "user_id missing in claims")
	// }
	// authorId := uint64(userId)

	article, err := s.uc.CreArticle(ctx, req.AuthorId, req.Title, req.Content)
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

func (s *ContentService) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleReply, error) {
	article, err := s.uc.GetArticle(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetArticleReply{
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

func (s *ContentService) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleReply, error) {
	article, err := s.uc.UpdateArticle(ctx, req.Id, req.AuthorId, req.Title, req.Content)
	if err != nil {
		return nil, err

	}
	return &pb.UpdateArticleReply{
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

func (s *ContentService) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticlReply, error) {
	err := s.uc.DeleteArticle(ctx, req.Id, req.AuthorId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteArticlReply{}, nil
}
