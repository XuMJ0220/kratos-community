package biz

import (
	"context"
	v1 "kratos-community/api/user/v1"
	"kratos-community/internal/conf"
	"kratos-community/internal/pkg/bcrypt"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

// 要有一个Usercase
type mockUserRepo struct {
}

// 预先计算好的哈希密码，避免每次重新生成
var (
	hashedPasswordCorrect, _ = bcrypt.HashPassword("correct_password")
)

func (r *mockUserRepo) GetUserByUserName(ctx context.Context, userName string) (*GetUserResponse, error) {
	// 模拟“用户不存在场景”
	if userName == "not_exist_user" {
		return nil, ErrUserNotFound // 直接返回业务错误
	}
	// 模拟“用户存在”场景
	if userName == "exist_user" {
		return &GetUserResponse{
			Password: hashedPasswordCorrect,
			UserInfo: v1.UserInfo{
				Id:       1,
				UserName: "exist_user",
				Email:    "<EMAIL>",
			},
		}, nil
	}
	return nil, nil
}

// 这里是不需要实现的，但是必须写上，因为我们是需要实现接口
func (r *mockUserRepo) CreateUser(ctx context.Context, ru *RegisterUser) error {
	return nil
}

// TestLogin
func TestLogin(t *testing.T) {
	// 创建一个UserUsecase实例，并且注入我们手写的mockRepo
	mockAuth := &conf.Auth{
		JwtSecret: "test_jwt_secret",
	}

	uc := NewUserUsecase(&mockUserRepo{}, log.DefaultLogger, mockAuth)
	// 表格驱动测试
	testCase := []struct {
		name        string // 测试用例名称
		userName    string
		password    string
		expectError bool // 是否期望得到一个错误
	}{
		{
			name:        "用户不存在",
			userName:    "not_exist_user",
			password:    "any_password",
			expectError: true,
		},
		{
			name:        "密码错误",
			userName:    "exist_user",
			password:    "wrong_password",
			expectError: true,
		},
		{
			name:        "登录成功",
			userName:    "exist_user",
			password:    "correct_password",
			expectError: false,
		},
	}

	// 遍历并运行所有测试用例
	for _, tc := range testCase {
		// 创建一个局部变量副本，防止闭包捕获问题
		tc := tc 
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // 加上这行可以让测试并行，速度更快
			// 调用我们测试的Login方法
			_, err := uc.Login(context.Background(), tc.userName, tc.password)
			// 断言检查结果是否符合我们的预期
			if (err != nil) != tc.expectError {
				t.Errorf("Login() error = %v, expectError %v", err, tc.expectError)
			}
		})
	}
}
