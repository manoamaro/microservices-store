package helpers

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/manoamaro/microservices-store/auth_service/models"
	"os"
	"strconv"
	"time"
)

func getJWTSecret() []byte {
	if value, exists := os.LookupEnv("JWT_SECRET"); exists {
		return []byte(value)
	}
	return []byte("My Secret")
}

func getJWTSecretFunc(_ *jwt.Token) (interface{}, error) {
	return getJWTSecret(), nil
}

func GetClaimsFromToken(rawToken string) (*models.UserClaims, error) {
	if token, err := jwt.ParseWithClaims(rawToken, &models.UserClaims{}, getJWTSecretFunc); err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	} else if userValues := token.Claims.(*models.UserClaims); userValues == nil {
		return nil, errors.New("invalid payload")
	} else {
		return userValues, nil
	}
}

func GetClaimsFromRefreshToken(rawToken string) (*jwt.RegisteredClaims, error) {
	if token, err := jwt.ParseWithClaims(rawToken, &jwt.RegisteredClaims{}, getJWTSecretFunc); err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	} else if claims := token.Claims.(*jwt.RegisteredClaims); claims == nil {
		return nil, errors.New("invalid payload")
	} else {
		return claims, nil
	}
}

func CreateTokens(authId uint, audience []string, flags []string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.Itoa(int(authId)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Audience:  audience,
		},
		AuthInfo: models.AuthInfo{
			Flags: flags,
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%d", authId),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
	})

	if accessTokenSigned, err := accessToken.SignedString(getJWTSecret()); err != nil {
		return "", "", err
	} else if refreshTokenSigned, err := refreshToken.SignedString(getJWTSecret()); err != nil {
		return "", "", err
	} else {
		return accessTokenSigned, refreshTokenSigned, nil
	}
}
