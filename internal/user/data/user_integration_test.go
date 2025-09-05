package data

import (
	"context"
	pb "kratos-community/api/user/v1"
	"kratos-community/internal/conf"
	"kratos-community/internal/pkg/testcontainerstools"
	"kratos-community/internal/user/biz"
	"kratos-community/internal/user/service"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUserIntegration(t *testing.T) {

	ctx := context.Background()
	// 启动一个临时的MySQL容器
	db, cleanup, err := testcontainerstools.SetupTestDB[User](ctx)
	// 检查错误
	assert.NoError(t, err, "setupTestDB should not return an error")
	// 最后清理容器
	defer cleanup()
	// 创建一个Data，连接了临时数据库
	d := &Data{
		db1: db,
	}
	// 创建一个UserRepo
	r := NewUserRepo(d, log.DefaultLogger)
	// 创建一个UserUsecase
	uc := biz.NewUserUsecase(r, log.DefaultLogger, &conf.Auth{JwtSecret: "test_secret"})
	// 创建一个UserService
	service := service.NewUserService(uc)
	// 准备测试数据
	registerReq := &pb.RegisterUserRequest{
		Email:      "test@example.com",
		Password:   "test_password",
		RePassword: "test_password",
		UserName:   "test_user",
	}
	// 调用我们要测试的业务逻辑
	_, err = service.RegisterUser(ctx, registerReq)
	// 检查错误
	assert.NoError(t, err, "RegisterUser should not return an error")
	// 直接查数据库，验证是否真的被写入了
	var User User
	result := db.Where("user_name = ?", "test_user").First(&User)

	assert.NoError(t, result.Error, "should find the saved user in db")
	assert.Equal(t, "test_user", User.UserName)
	assert.NotEmpty(t, User.Password, "password should be hashed and not empty")

}
