package repositories

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/models"
	"github.com/manoamaro/microservices-store/commons/pkg/collections"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepository interface {
	GetTokenFromRequest(r *http.Request) (*models.UserClaims, string, error)
	GetClaimsFromToken(rawToken string) (*models.UserClaims, error)
	GetClaimsFromRefreshToken(rawToken string) (*jwt.RegisteredClaims, error)
	CreateAuth(email string, plainPassword string) (auth *models.Auth, err error)
	Authenticate(email string, plainPassword string) (auth *models.Auth, found bool)
	InvalidateToken(token *models.UserClaims, rawToken string) error
	IsInvalidatedToken(rawToken string) bool
	CreateTokens(authId uint) (accessToken, refreshToken string, err error)
}

type DefaultAuthRepository struct {
	context     context.Context
	redisClient *redis.Client
	db          *sql.DB
	ormDB       *gorm.DB
}

func NewDefaultAuthRepository(db *sql.DB, redisClient *redis.Client) AuthRepository {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return &DefaultAuthRepository{
		context:     context.Background(),
		redisClient: redisClient,
		db:          db,
		ormDB:       gormDB,
	}
}

func (s *DefaultAuthRepository) CreateAuth(email string, plainPassword string) (auth *models.Auth, err error) {
	s.ormDB.Preload(clause.Associations).Where(&models.Auth{Email: email}).First(&auth)
	if auth != nil && auth.ID > 0 {
		return nil, errors.New("user already exists")
	}
	salt := strconv.FormatInt(time.Now().UnixNano(), 16)
	auth = &models.Auth{
		Email:    email,
		Password: CalculatePasswordHash(plainPassword, salt),
		Salt:     salt,
	}

	result := s.ormDB.Create(auth)
	if result.Error != nil {
		return nil, result.Error
	}
	return auth, nil
}

func (s *DefaultAuthRepository) Authenticate(email string, plainPassword string) (auth *models.Auth, found bool) {
	s.ormDB.Preload(clause.Associations).Where(&models.Auth{Email: email}).First(&auth)
	if auth.ID == 0 {
		return nil, false
	} else if passwordHash := CalculatePasswordHash(plainPassword, auth.Salt); passwordHash != auth.Password {
		return nil, false
	} else {
		return auth, true
	}
}

func (s *DefaultAuthRepository) InvalidateToken(token *models.UserClaims, rawToken string) error {
	key := s.getRedisInvalidTokenKey(rawToken)
	expiration := time.Now().Sub(token.ExpiresAt.Time)
	err := s.redisClient.Set(s.context, key, true, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultAuthRepository) IsInvalidatedToken(rawToken string) bool {
	key := s.getRedisInvalidTokenKey(rawToken)
	val, _ := s.redisClient.Get(s.context, key).Bool()
	return val
}

func (s *DefaultAuthRepository) getRedisInvalidTokenKey(rawToken string) string {
	return fmt.Sprintf("token.invalid.%s", rawToken)
}

func (s *DefaultAuthRepository) CreateTokens(authId uint) (string, string, error) {
	var auth models.Auth
	s.ormDB.Preload(clause.Associations).Where("id = ?", authId).First(&auth)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.Itoa(int(authId)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			Audience:  collections.MapTo(auth.Domains, func(audience models.Domain) string { return audience.Domain }),
		},
		AuthInfo: models.AuthInfo{
			Flags: collections.MapTo(auth.Flags, func(flag models.Flag) string { return flag.Name }),
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        strconv.Itoa(int(auth.ID)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
	})

	if accessTokenSigned, err := accessToken.SignedString(helpers.GetJWTSecret()); err != nil {
		return "", "", err
	} else if refreshTokenSigned, err := refreshToken.SignedString(helpers.GetJWTSecret()); err != nil {
		return "", "", err
	} else {
		return accessTokenSigned, refreshTokenSigned, nil
	}
}

func (s *DefaultAuthRepository) GetTokenFromRequest(r *http.Request) (*models.UserClaims, string, error) {
	if rawToken, err := request.AuthorizationHeaderExtractor.ExtractToken(r); err != nil {
		return nil, "", err
	} else if userClaims, err := s.GetClaimsFromToken(rawToken); err != nil {
		return nil, "", err
	} else {
		return userClaims, rawToken, nil
	}
}

func (s *DefaultAuthRepository) GetClaimsFromToken(rawToken string) (*models.UserClaims, error) {
	if token, err := jwt.ParseWithClaims(rawToken, &models.UserClaims{}, helpers.GetJWTSecretFunc); err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	} else if userValues := token.Claims.(*models.UserClaims); userValues == nil {
		return nil, errors.New("invalid payload")
	} else {
		return userValues, nil
	}
}

func (s *DefaultAuthRepository) GetClaimsFromRefreshToken(rawToken string) (*jwt.RegisteredClaims, error) {
	if token, err := jwt.ParseWithClaims(rawToken, &jwt.RegisteredClaims{}, helpers.GetJWTSecretFunc); err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	} else if claims := token.Claims.(*jwt.RegisteredClaims); claims == nil {
		return nil, errors.New("invalid payload")
	} else {
		return claims, nil
	}
}

func CalculatePasswordHash(plainPassword string, salt string) string {
	h := sha256.New()
	h.Write([]byte(plainPassword))
	h.Write([]byte(salt))
	r := fmt.Sprintf("%x", h.Sum(nil))
	return r
}
