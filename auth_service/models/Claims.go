package models

import "github.com/golang-jwt/jwt/v4"

type AuthInfo struct {
	Flags []string
}

type UserClaims struct {
	jwt.RegisteredClaims
	AuthInfo
}
