package data

import (
    "time"
    "gorm.io/gorm"
)

type Article struct {
    ID        uint64         `gorm:"primaryKey;autoIncrement;comment:文章ID" json:"id"`
    AuthorID  uint64         `gorm:"not null;index:idx_author_id;comment:作者的用户ID" json:"author_id"`
    Title     string         `gorm:"type:varchar(100);not null;comment:文章标题" json:"title"` // collate 可以在表级别设置，字段上非必须
    Content   string         `gorm:"type:longtext;not null;comment:文章内容" json:"content"`
    CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"` // GORM 会自动处理创建时间
    UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"` // GORM 会自动处理更新时间
    DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:删除时间" json:"deleted_at"` // 修正：补充索引标签
}

// OutboxMessage 事务性发件箱表
type OutboxMessage struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;comment:主键"`
	Topic        string    `gorm:"type:varchar(255);not null;comment:kafka topic"`
	MessageKey   string    `gorm:"type:varchar(255);not null;comment:消息 Key"`
	MessageValue string    `gorm:"type:varchar(255);not null;comment:消息 Value"`
	Status       uint8     `gorm:"type:tinyint;not null;default:0;comment:状态: 0-待发送, 1-已发送"`
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

// TableName 指定表名
func (Article) TableName() string {
    return "articles"
}

// TableName 指定 GORM 应该将此模型映射到哪个表。
func (OutboxMessage) TableName() string {
	return "outbox_messages"
}
