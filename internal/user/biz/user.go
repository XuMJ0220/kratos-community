package biz

import (
	"context"
	v1 "kratos-community/api/user/v1"
	"kratos-community/internal/conf"
	"kratos-community/internal/pkg/bcrypt"
	myjwt "kratos-community/internal/pkg/jwt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
)

// 定义一些错误
var (
	ErrUserPasswordNotMatch = errors.BadRequest("User_Bad_Request", "密码不一致")
	ErrUserNotFound         = errors.NotFound("User_Not_Found", "用户名或密码错误")
	ErrUserAlreadyExist     = errors.Conflict("User_Already_Exist", "用户或邮箱已存在")
	ErrGenerateToken        = errors.InternalServer("GenerateToken_Error", "生成Token失败")
)

// UserRepo 与数据库交互的接口
type UserRepo interface {
	CreateUser(ctx context.Context, ru *RegisterUser) error
	GetUserByUserName(ctx context.Context, userName string) (*GetUserResponse, error)
}

type RegisterUser struct {
	UserName   string
	Password   string
	Email      string
	RePassword string
}

type UserUsecase struct {
	repo      UserRepo
	log       *log.Helper
	jwtSecret string
}

// LoginResponse 登陆业务逻辑响应
type LoginResponse struct {
	Token    string
	UserInfo *v1.UserInfo
}

type GetUserResponse struct {
	UserInfo v1.UserInfo
	Password string
}

func NewRegisterUser(userName, password, email, rePassword string) *RegisterUser {
	return &RegisterUser{UserName: userName, Password: password, Email: email, RePassword: rePassword}
}

func NewUserUsecase(repo UserRepo, logger log.Logger, jwtSecret *conf.Auth) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger), jwtSecret: jwtSecret.JwtSecret}
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

func (uc *UserUsecase) Login(ctx context.Context, userName, password string) (*LoginResponse, error) {
	// 1. 根据用户名数据库查询用户信息
	user, err := uc.repo.GetUserByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	// 2. 判断密码是否正确
	if !bcrypt.CheckPasswordHash(password, user.Password) {
		return nil, ErrUserNotFound
	}
	// 3. 生成Token
	claims := myjwt.NewCustomClaims(user.UserInfo.Id, user.UserInfo.UserName, 24*time.Hour)
	token, err := myjwt.GenerateToken(claims, uc.jwtSecret, jwt.SigningMethodHS256)
	if err != nil {
		uc.log.Errorf("Login: GenerateToken error : %v", err)
		return nil, ErrGenerateToken
	}
	return &LoginResponse{
		Token:    token,
		UserInfo: &user.UserInfo,
	}, nil
}
