package data

import "time"

// Like 对应于数据库中的 `likes` 表
type Like struct {
	// 点赞的用户ID, 是复合主键的一部分
	UserID    uint64    `gorm:"column:user_id;type:bigint unsigned;primaryKey"`
	
	// 被点赞的文章ID, 是复合主键的一部分, 同时也是一个索引
	ArticleID uint64    `gorm:"column:article_id;type:bigint unsigned;primaryKey;index:idx_article_id"`
	
	// 点赞时间, GORM 会在创建时自动填充
	CreatedAt time.Time `gorm:"column:created_at"`
}

// TableName 方法指定了 Like 结构体对应的数据库表名
func (l *Like) TableName() string {
	return "likes"
}