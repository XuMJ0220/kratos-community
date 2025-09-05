package data

import "time"

// Relation 对应于数据库中的 `relations` 表
type Relation struct {
	// 粉丝ID (发起关注的用户), 是复合主键的一部分
	FollowerID uint64 `gorm:"column:follower_id;type:bigint unsigned;primaryKey"`

	// 被关注者ID, 是复合主键的一部分, 同时也是一个独立的索引
	FollowingID uint64 `gorm:"column:following_id;type:bigint unsigned;primaryKey;index:idx_following_id"`

	// 关注时间, GORM 会在创建记录时自动填充
	CreatedAt time.Time `gorm:"column:created_at"`
}

// TableName 方法指定了 Relation 结构体对应的数据库表名
func (r *Relation) TableName() string {
	return "relations"
}