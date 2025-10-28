package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "secret" // fallback khi chưa có biến môi trường
	}
	jwtSecret = []byte(s)
}

// Custom error để thay cho ErrTokenInvalid/ErrTokenInvalidMethod cũ
var (
	ErrInvalidTokenMethod = errors.New("invalid token signing method")
	ErrInvalidToken       = errors.New("invalid token")
)

// GenerateToken tạo JWT chứa userID trong claims
func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken kiểm tra token và trả về userID nếu hợp lệ
func ParseToken(tokenStr string) (uint, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Kiểm tra phương thức ký
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenMethod
		}
		return jwtSecret, nil
	})
	if err != nil {
		// Nếu là lỗi chữ ký hoặc token hết hạn, trả về lỗi gốc
		return 0, err
	}

	if !tok.Valid {
		return 0, ErrInvalidToken
	}

	// Lấy userID từ claims
	if claims, ok := tok.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(float64); ok {
			return uint(sub), nil
		}
	}

	return 0, ErrInvalidToken
}
