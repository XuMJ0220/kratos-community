package data

import (
	"kratos-community/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewInteractionRepo,NewDB,NewRedisClient)

// Data .
type Data struct {
	// TODO wrapped database client
	db1  *gorm.DB      // Mysql
	rdb1 *redis.Client // Redis
}

// NewData .
func NewData(db *gorm.DB,rdb *redis.Client,logger log.Logger) (*Data, func(), error) {

	// dsn := c.Databases["user_1"].Source

	// // 连接数据库
	// db1, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.NewHelper(logger).Errorf("failed to connect database: %v", err)
	// 	return nil, nil, err
	// }
	// log.NewHelper(logger).Infof("数据库连接成功")

	logHelper:=log.NewHelper(logger)
	d:=&Data{
		db1: db,
		rdb1: rdb,
	}

	// 定义一个 cleanup 函数，用于关闭数据库连接
	cleanup := func() {
		logHelper.Info("closing the data resources")
		// 关闭数据库连接
		sqlDB, _ := d.db1.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		// 关闭Redis连接
		if err:=d.rdb1.Close();err!=nil{
			logHelper.Errorf("failed to close redis: %v", err)
		}
	}

	return d, cleanup, nil
}

// NewBD 创建Mysql客户端
func NewDB(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
	logHelper := log.NewHelper(logger)
	db, err := gorm.Open(mysql.Open(c.Databases["user_1"].Source), &gorm.Config{})
	if err != nil {
		logHelper.Errorf("failed to connect database: %v", err)
		return nil, err
	}
	logHelper.Info("数据库连接成功")
	return db, nil
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(c *conf.Data, logger log.Logger) (*redis.Client, error) {
	logHelper := log.NewHelper(logger)
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})
	// 测试连接
	if err := rdb.Ping().Err(); err != nil {
		logHelper.Errorf("failed to connect redis: %v", err)
		return nil, err
	}
	logHelper.Info("redis连接成功")
	return rdb, nil
}
