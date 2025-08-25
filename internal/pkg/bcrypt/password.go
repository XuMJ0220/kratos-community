package bcrypt

// internal/pkg/utils/password.go

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 对密码进行哈希处理
// password: 明文密码
// 返回值: 哈希后的密码字符串 和 错误
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword 会自动生成盐并将其包含在哈希值中
	// bcrypt.DefaultCost 是默认的计算成本，数值越高，计算越慢，也就越安全
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash 验证明文密码和哈希值是否匹配
// password: 用户输入的明文密码
// hash: 数据库中存储的哈希密码
// 返回值: true 表示匹配成功，false 表示失败
func CheckPasswordHash(password, hash string) bool {
	// bcrypt.CompareHashAndPassword 会从 hash 中提取出盐，
	// 然后用相同的盐对 password 进行哈希，最后比较两个哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// 如果 err 为 nil，表示密码匹配成功
	return err == nil
}