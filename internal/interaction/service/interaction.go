package service

import (
	"context"

	pb "kratos-community/api/interaction/v1"
	"kratos-community/internal/interaction/biz"
)

type InteractionService struct {
	pb.UnimplementedInteractionServer

	uc *biz.InteractionUsecase
}

func NewInteractionService(uc *biz.InteractionUsecase) *InteractionService {
	return &InteractionService{uc: uc}
}

func (s *InteractionService) LikeArticle(ctx context.Context, req *pb.LikeArticleRequest) (*pb.LikeArticleReply, error) {
	err:=s.uc.LikeArticle(ctx,req.UserId,req.Id)
	if err!=nil{
		return nil,err
	}
	return &pb.LikeArticleReply{}, nil
}
func (s *InteractionService) UnlikeArticle(ctx context.Context, req *pb.UnlikeArticleRequest) (*pb.UnlikeArticleReply, error) {
	err:=s.uc.UnLikeArticle(ctx,req.UserId,req.Id)
	if err!=nil{
		return nil,err
	}
	return &pb.UnlikeArticleReply{}, nil
}
