package data

import (
	"context"
	"kratos-community/internal/user/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-sql-driver/mysql"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateUser 创建用户
func (r *userRepo) CreateUser(ctx context.Context, ru *biz.RegisterUser) error {

	// 创建一个 User 模型
	user := User{
		UserName: ru.UserName,
		Password: ru.Password,
		Email:    ru.Email,
	}
	// 创建用户
	result := r.data.db1.WithContext(ctx).Create(&user)
	if result.Error != nil {

		var mysqlErr *mysql.MySQLError

		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			return biz.ErrUserAlreadyExist
		}
	}
	// 创建成功
	return nil
}
