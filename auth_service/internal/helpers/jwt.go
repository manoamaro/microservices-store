package helpers

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/manoamaro/microservices-store/auth_service/models"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"strconv"
	"time"
)

var ErrJWTInvalidPayload = errors.New("invalid payload")
var ErrJWTInvalidToken = errors.New("invalid token")

func getJWTSecret() []byte {
	secretStr := helpers.GetEnv("JWT_SECRET", "My Secret")
	return []byte(secretStr)
}

func getJWTSecretFunc(_ *jwt.Token) (interface{}, error) {
	return getJWTSecret(), nil
}

func GetClaimsFromToken(rawToken string) (*models.UserClaims, error) {
	if token, err := jwt.ParseWithClaims(rawToken, &models.UserClaims{}, getJWTSecretFunc); err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, ErrJWTInvalidToken
	} else if userValues := token.Claims.(*models.UserClaims); userValues == nil {

		return nil, ErrJWTInvalidPayload
	} else {
		return userValues, nil
	}
}

func GetClaimsFromRefreshToken(rawToken string) (*jwt.RegisteredClaims, error) {
	if token, err := jwt.ParseWithClaims(rawToken, &jwt.RegisteredClaims{}, getJWTSecretFunc); err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, ErrJWTInvalidToken
	} else if claims := token.Claims.(*jwt.RegisteredClaims); claims == nil {
		return nil, ErrJWTInvalidPayload
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
