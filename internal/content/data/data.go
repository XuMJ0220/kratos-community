package data

import (
	"kratos-community/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData,NewContentRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db1 *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {

	dsn := c.Databases["user_1"].Source

	// 连接数据库
	db1, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.NewHelper(logger).Errorf("failed to connect database: %v", err)
		return nil, nil, err
	}
	log.NewHelper(logger).Infof("user服务连接数据库user_1成功")
	// 定义一个 cleanup 函数，用于关闭数据库连接
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")

		sqlDB, _ := db1.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	// 创建Data实例
	d := &Data{db1: db1}
	return d, cleanup, nil
}
