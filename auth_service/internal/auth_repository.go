package internal

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	models2 "manoamaro.github.com/auth_service/internal/models"
	"manoamaro.github.com/commons"
	"net/http"
	"strconv"
	"time"
)

type AuthRepository interface {
	GetTokenFromRequest(r *http.Request) (*UserClaims, string, error)
	CreateAuth(email string, plainPassword string) (auth *models2.Auth, err error)
	Authenticate(email string, plainPassword string) (auth *models2.Auth, found bool)
	InvalidateToken(token *UserClaims, rawToken string) error
	CheckToken(rawToken string) bool
	CreateToken(auth *models2.Auth) (string, error)
}

type DefaultAuthRepository struct {
	context     context.Context
	redisClient *redis.Client
	db          *sql.DB
	ormDB       *gorm.DB
}

func NewDefaultAuthRepository() AuthRepository {
	db, err := sql.Open("postgres", commons.GetEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err := gormDB.AutoMigrate(&models2.Auth{}, &models2.Flag{}, &models2.Audience{}, &models2.Domain{}); err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     commons.GetEnv("REDIS_URL", "localhost:6379"),
		Username: commons.GetEnv("REDIS_USERNAME", ""),
		Password: commons.GetEnv("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	return &DefaultAuthRepository{
		context:     context.Background(),
		redisClient: redisClient,
		db:          db,
		ormDB:       gormDB,
	}
}

func (s *DefaultAuthRepository) CreateAuth(email string, plainPassword string) (auth *models2.Auth, err error) {
	s.ormDB.Preload(clause.Associations).Where(&models2.Auth{Email: email}).First(&auth)
	if auth != nil && auth.ID > 0 {
		return nil, errors.New("user already exists")
	}
	salt := strconv.FormatInt(time.Now().UnixNano(), 16)
	auth = &models2.Auth{
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

func (s *DefaultAuthRepository) Authenticate(email string, plainPassword string) (auth *models2.Auth, found bool) {
	s.ormDB.Preload(clause.Associations).Where(&models2.Auth{Email: email}).First(&auth)
	if auth.ID == 0 {
		return nil, false
	} else if passwordHash := CalculatePasswordHash(plainPassword, auth.Salt); passwordHash != auth.Password {
		return nil, false
	} else {
		return auth, true
	}
}

func (s *DefaultAuthRepository) InvalidateToken(token *UserClaims, rawToken string) error {
	key := s.getRedisInvalidTokenKey(rawToken)
	expiration := time.Now().Sub(token.ExpiresAt.Time)
	err := s.redisClient.Set(s.context, key, true, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultAuthRepository) CheckToken(rawToken string) bool {
	key := s.getRedisInvalidTokenKey(rawToken)
	val, _ := s.redisClient.Get(s.context, key).Bool()
	return val
}

func (s *DefaultAuthRepository) getRedisInvalidTokenKey(rawToken string) string {
	return fmt.Sprintf("token.invalid.%s", rawToken)
}

func (s *DefaultAuthRepository) CreateToken(auth *models2.Auth) (string, error) {
	var audiences []models2.Audience
	s.ormDB.Preload(clause.Associations).Where(&models2.Audience{Auth: auth}).First(&audiences)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		jwt.RegisteredClaims{
			ID:        strconv.Itoa(int(auth.ID)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			Audience:  commons.MapTo(audiences, func(audience models2.Audience) string { return audience.Domain.Domain }),
		},
		AuthInfo{
			Flags: commons.MapTo(auth.Flags, func(flag models2.Flag) string { return flag.Name }),
		},
	})
	return token.SignedString(GetJWTSecret())
}

func (s *DefaultAuthRepository) GetTokenFromRequest(r *http.Request) (*UserClaims, string, error) {
	rawToken, err := request.AuthorizationHeaderExtractor.ExtractToken(r)
	if err != nil {
		return nil, "", err
	}

	token, err := jwt.ParseWithClaims(rawToken, &UserClaims{}, GetJWTSecretFunc)
	if err != nil {
		return nil, "", err
	}

	if !token.Valid {
		return nil, "", errors.New("invalid token")
	}

	invalidated := s.CheckToken(rawToken)
	if invalidated {
		return nil, "", errors.New("token invalidated")
	}

	userValues := token.Claims.(*UserClaims)
	if userValues == nil {
		return nil, "", errors.New("invalid payload")
	}

	return userValues, rawToken, nil
}

func CalculatePasswordHash(plainPassword string, salt string) string {
	h := sha256.New()
	h.Write([]byte(plainPassword))
	h.Write([]byte(salt))
	r := fmt.Sprintf("%x", h.Sum(nil))
	return r
}
