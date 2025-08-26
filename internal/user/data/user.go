package data

import (
	"context"
	v1 "kratos-community/api/user/v1"
	"kratos-community/internal/user/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
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

func (r *userRepo) GetUserByUserName(ctx context.Context, userName string) (*biz.GetUserResponse, error) {

	// 查询
	user, err := gorm.G[User](r.data.db1).First(ctx)
	
	if err!=nil{
		r.log.Errorf("GetUserByUserName: %v", err) // 输出错误日志
		if errors.Is(err, gorm.ErrRecordNotFound) { // 不存在该行数据
			return nil, biz.ErrUserNotFound
		}else{ // 其他错误
			return nil,err 
		}
	}

	return &biz.GetUserResponse{
		UserInfo:v1.UserInfo{
			Id:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
		},
		Password: user.Password,
	}, nil
}
