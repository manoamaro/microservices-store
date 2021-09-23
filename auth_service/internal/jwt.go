package internal

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

func GetJWTSecret() []byte {
	if value, exists := os.LookupEnv("JWT_SECRET"); exists {
		return []byte(value)
	}
	return []byte("My Secret")
}

func GetTokenSigned(userId string, userEmail string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		&jwt.StandardClaims{
			Id:        userId,
			ExpiresAt: time.Now().Add(time.Hour * 168).Unix(),
		},
		UserInfo{
			Email:  userEmail,
			Access: []string{""},
		},
	})
	return token.SignedString(GetJWTSecret())
}

type UserInfo struct {
	Email  string
	Access []string
}

type UserClaims struct {
	*jwt.StandardClaims
	UserInfo
}
