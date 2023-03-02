package internal

import (
	"database/sql"
	"fmt"
	"github.com/manoamaro/microservices-store/auth_service/internal/controllers"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/golang-migrate/migrate/v4"
	"github.com/manoamaro/microservices-store/auth_service/internal/repositories"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
)

type Application struct {
	db             *sql.DB
	redisClient    *redis.Client
	r              *gin.Engine
	authRepository repositories.AuthRepository
	authController *controllers.AuthController
}

func NewApplication() *Application {
	postgresUrl := helpers.GetEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/auth_service?sslmode=disable")

	db, err := sql.Open("postgres", postgresUrl)

	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     helpers.GetEnv("REDIS_URL", "localhost:6379"),
		Username: helpers.GetEnv("REDIS_USERNAME", ""),
		Password: helpers.GetEnv("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	r := gin.Default()

	authRepository := repositories.NewDefaultAuthRepository(db, redisClient)
	authController := controllers.NewAuthController(r, authRepository)

	return &Application{
		db:             db,
		redisClient:    redisClient,
		r:              r,
		authRepository: authRepository,
		authController: authController,
	}
}

func (a *Application) RunMigrations() {
	migration, err := NewMigration(a.db)
	if err != nil {
		log.Fatal(err)
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func (a *Application) Run(c chan error) {
	port := helpers.GetEnv("PORT", "8080")

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			c <- err
		}
	}()
}
