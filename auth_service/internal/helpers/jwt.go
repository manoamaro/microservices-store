package helpers

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
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
