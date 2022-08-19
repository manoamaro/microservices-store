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

func GetJWTSecretFunc(_ *jwt.Token) (interface{}, error) {
	return GetJWTSecret(), nil
}

func GetTokenSigned(authId string, roles []string, flags []string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		jwt.RegisteredClaims{
			ID:        authId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
		AuthInfo{
			Roles: roles,
			Flags: flags,
		},
	})
	return token.SignedString(GetJWTSecret())
}

type AuthInfo struct {
	Roles []string
	Flags []string
}

type UserClaims struct {
	jwt.RegisteredClaims
	AuthInfo
}
