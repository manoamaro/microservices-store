package internal

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"manoamaro.github.com/auth_service/models"
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

	if err := gormDB.AutoMigrate(&models.Flag{}, &models.Role{}, &models.Auth{}); err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     GetENV("REDS_URL", "localhost:6379"),
		Username: GetENV("REDS_USERNAME", ""),
		Password: GetENV("REDS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	return &AuthService{
		context:     context.Background(),
		redisClient: redisClient,
		db:          db,
		ormDB:       gormDB,
	}
}

func (s *AuthService) FindAuth(userId uint64) (auth models.Auth, found bool) {
	s.ormDB.Preload(clause.Associations).Where(&models.Auth{UserId: userId}).First(&auth)
	if auth.ID == 0 {
		return auth, false
	}
	return auth, true
}
