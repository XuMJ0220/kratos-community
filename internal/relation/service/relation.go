package service

import (
	"context"

	pb "kratos-community/api/relation/v1"
	"kratos-community/internal/relation/biz"
)

type RelationService struct {
	pb.UnimplementedRelationServer
	uc *biz.RelationUsecase
}

func NewRelationService(uc *biz.RelationUsecase) *RelationService {
	return &RelationService{uc: uc}
}

func (s *RelationService) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserReply, error) {
	err:=s.uc.FollowUser(ctx,req.UserId,req.Id)
	if err!=nil{
		return nil,err
	}
	return &pb.FollowUserReply{}, nil
}
func (s *RelationService) UnfollowUser(ctx context.Context, req *pb.UnfollowUserRequest) (*pb.UnfollowUserReply, error) {
	err:=s.uc.UnfollowUser(ctx,req.UserId,req.Id)
	if err!=nil{
		return nil,err
	}
	return &pb.UnfollowUserReply{}, nil
}
func (s *RelationService) ListFollowings(ctx context.Context, req *pb.ListFollowingsRequest) (*pb.ListFollowingsReply, error) {
	return &pb.ListFollowingsReply{}, nil
}
func (s *RelationService) ListFollowers(ctx context.Context, req *pb.ListFollowersRequest) (*pb.ListFollowersReply, error) {
	return &pb.ListFollowersReply{}, nil
}
