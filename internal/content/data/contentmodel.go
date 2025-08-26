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

// TableName 指定表名
func (Article) TableName() string {
    return "articles"
}