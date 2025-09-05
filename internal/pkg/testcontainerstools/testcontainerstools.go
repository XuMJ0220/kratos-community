package testcontainerstools

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait" // 导入 wait 包
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupTestDB 启动一个临时的MySQL容器
func SetupTestDB[T any](ctx context.Context) (*gorm.DB, func(), error) {
	// 定义数据库连接信息
	dbName := "test_db"
	dbUser := "test_user"
	dbPassword := "test_password"

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root_password",
			"MYSQL_DATABASE":      dbName,
			"MYSQL_USER":          dbUser,
			"MYSQL_PASSWORD":      dbPassword,
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(time.Minute * 2),
	}

	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start mysql container: %w", err)
	}

	// 手动构建 DSN
	// 1. 获取容器动态映射的主机和端口
	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container host: %w", err)
	}
	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container port: %w", err)
	}     

	// 2. 使用 fmt.Sprintf 拼接成标准的 DSN 字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", dbUser, dbPassword, host, port.Port(), dbName)

	// 3. 连接到这个临时的数据库 (gormMysql 是 gorm.io/driver/mysql 的别名)
	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	// 执行数据库迁移（建表）
	// 使用泛型类型T来创建模型实例
	var model T
	err = db.AutoMigrate(&model)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	// 定义 cleanup 函数
	cleanup := func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Errorf("failed to terminate mysql container: %v", err)
		}
	}

	return db, cleanup, nil
}
