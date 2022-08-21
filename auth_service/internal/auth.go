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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"manoamaro.github.com/auth_service/models"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type AuthService struct {
	context     context.Context
	redisClient *redis.Client
	db          *sql.DB
	ormDB       *gorm.DB
}

func NewAuthService() *AuthService {
	db, err := sql.Open("postgres", GetENV("DB_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err := gormDB.AutoMigrate(&models.Auth{}, &models.Flag{}, &models.Audience{}, &models.Domain{}); err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     GetENV("REDIS_URL", "localhost:6379"),
		Username: GetENV("REDIS_USERNAME", ""),
		Password: GetENV("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	return &AuthService{
		context:     context.Background(),
		redisClient: redisClient,
		db:          db,
		ormDB:       gormDB,
	}
}

func (s *AuthService) Authenticate(email string, plainPassword string) (auth *models.Auth, found bool) {
	s.ormDB.Preload(clause.Associations).Where(&models.Auth{Email: email}).First(&auth)
	if auth.ID == 0 {
		return nil, false
	} else if passwordHash := CalculatePasswordHash(plainPassword, auth.Salt); passwordHash != auth.Password {
		return nil, false
	} else {
		return auth, true
	}
}

func (s *AuthService) CreateAuth(email string, plainPassword string) (auth *models.Auth, err error) {
	s.ormDB.Preload(clause.Associations).Where(&models.Auth{Email: email}).First(&auth)
	if auth != nil && auth.ID > 0 {
		return nil, err
	}
	salt := strconv.FormatInt(rand.Int63(), 16)
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

func (s *AuthService) CreateToken(auth *models.Auth) (string, error) {
	var audiences []models.Audience
	s.ormDB.Preload(clause.Associations).Where(&models.Audience{Auth: auth}).First(&audiences)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		jwt.RegisteredClaims{
			ID:        strconv.Itoa(int(auth.ID)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			Audience:  MapTo(audiences, func(audience models.Audience) string { return audience.Domain.Domain }),
		},
		AuthInfo{
			Flags: MapTo(auth.Flags, func(flag models.Flag) string { return flag.Name }),
		},
	})
	return token.SignedString(GetJWTSecret())
}

func (s *AuthService) GetTokenFromRequest(r *http.Request) (*UserClaims, error) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, GetJWTSecretFunc, request.WithClaims(&UserClaims{}))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token not valid")
	}

	userValues := token.Claims.(*UserClaims)
	if userValues == nil {
		return nil, errors.New("invalid payload")
	}

	return userValues, nil
}

func CalculatePasswordHash(plainPassword string, salt string) string {
	h := sha256.New()
	h.Write([]byte(plainPassword))
	h.Write([]byte(salt))
	r := fmt.Sprintf("%x", h.Sum(nil))
	return r
}
