package service

import (
	"context"

	pb "kratos-community/api/user/v1"
	"kratos-community/internal/user/biz"
)

type UserService struct {
	pb.UnimplementedUserServer

	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserReply, error) {

	ru := biz.NewRegisterUser(req.UserName, req.Password, req.Email, req.RePassword)
	if err := s.uc.RegisterUser(ctx, ru); err != nil {
		return nil, err
	}
	return &pb.RegisterUserReply{}, nil
}
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	r, err := s.uc.Login(ctx, req.UserName, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		Token:    r.Token,
		UserInfo: r.UserInfo,
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	users,err := s.uc.ListUsers(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &pb.ListUsersReply{
		Users: users,
	}, nil
}
