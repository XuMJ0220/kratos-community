package data

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint64         `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
    UserName  string         `gorm:"type:varchar(30);uniqueIndex:uk_user_name;not null;collate:utf8mb4_bin;comment:用户名" json:"user_name"`
    Password  string         `gorm:"type:varchar(255);not null;collate:utf8mb4_unicode_ci;comment:密码" json:"password"`
    Email     string         `gorm:"type:varchar(255);uniqueIndex:uk_email;not null;collate:utf8mb4_unicode_ci;comment:邮箱" json:"email"`
    CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
    UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:删除时间（用于软删除）" json:"deleted_at"`
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}