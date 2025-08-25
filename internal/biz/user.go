package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

type registerUser struct{
	UserName string
	Password string
	Email string
	RePassword string
}

func NewRegisterUser(userName, password, email, rePassword string) *registerUser {
	return &registerUser{UserName: userName, Password: password, Email: email, RePassword: rePassword}
}

type UserUsecase struct {
	// repo
	log *log.Helper
}

func NewUserUsecase(logger log.Logger) *UserUsecase {
	return &UserUsecase{log:log.NewHelper(logger)}
}

func (uc *UserUsecase) RegisterUser(ctx context.Context,ru *registerUser) error{
	if ru.Password!=ru.RePassword{
		return errors.New("密码不一致")
	}
	uc.log.Infof("RegisterUser: %v", ru)
	return nil
} 