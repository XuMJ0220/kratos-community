package biz

import (
	"context"
	v1 "kratos-community/api/user/v1"
	"kratos-community/internal/pkg/bcrypt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// 定义一些错误
var (
	ErrUserPasswordNotMatch = errors.BadRequest("User_Bad_Request", "密码不一致")
	ErrUserNotFound         = errors.NotFound("User_Not_Found", "用户不存在")
	ErrUserAlreadyExist     = errors.Conflict("User_Already_Exist", "用户或邮箱已存在")
)

// UserRepo 与数据库交互的接口
type UserRepo interface {
	CreateUser(ctx context.Context, ru *RegisterUser) error
}

type RegisterUser struct {
	UserName   string
	Password   string
	Email      string
	RePassword string
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// LoginResponse 登陆业务逻辑响应
type LoginResponse struct {
	Token    string
	UserInfo *v1.UserInfo
}

func NewRegisterUser(userName, password, email, rePassword string) *RegisterUser {
	return &RegisterUser{UserName: userName, Password: password, Email: email, RePassword: rePassword}
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, ru *RegisterUser) error {
	if ru.Password != ru.RePassword {
		uc.log.Errorf("RegisterUser: Password : %v != RePassword : %v", ru.Password, ru.RePassword)
		return ErrUserPasswordNotMatch
	}
	uc.log.Infof("RegisterUser: %v", ru)

	// 在biz层用bcrypt对密码进行加密
	ru.Password, _ = bcrypt.HashPassword(ru.Password)
	err := uc.repo.CreateUser(ctx, ru) // 创建用户
	if err != nil {
		uc.log.Errorf("RegisterUser: CreateUser error : %v", err)
		return err
	}

	return nil
}

func (uc *UserUsecase) Login(ctx context.Context, userName, password string) error {
	return nil
}
