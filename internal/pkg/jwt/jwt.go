package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type customClaims struct {
	UserId               uint64
	UserName             string
	jwt.RegisteredClaims // 内嵌注册声明
}

// NewCustomCliams 创建自定义声明
// expireTime为过期时间
func NewCustomClaims(userId uint64, userName string, expireTime time.Duration) *customClaims {
	return &customClaims{
		UserId:   userId,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "kratos-community",
			Subject:   "user",
		},
	}
}

// GeneratoToken 生成 Token
// claims: 自定义声明
// secret: 密钥
// method: 签名方法
func GenerateToken( claims *customClaims, secret string, method jwt.SigningMethod) (string, error){
	token := jwt.NewWithClaims(method, claims) // 使用指定的签名方法和声明创建一个新的 Token
	ss,err:=token.SignedString([]byte(secret)) // 使用密钥签名并获取完整的编码后的字符串 Token
	if err!=nil{ 
		return "",err
	}
	return ss,nil
}

// ParseToken 解析 Token
// tokenString: Token 字符串
// secret: 密钥
func ParseToken(tokenString string, secret string) (*customClaims, error){
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrHashUnavailable
}