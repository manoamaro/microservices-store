package internal

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/golang-migrate/migrate/v4"
	"github.com/manoamaro/microservices-store/auth_service/internal/controllers"
	"github.com/manoamaro/microservices-store/auth_service/internal/repositories"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"log"
	"net/http"
	"time"
)

type Application struct {
	db             *sql.DB
	redisClient    *redis.Client
	r              *gin.Engine
	authRepository repositories.AuthRepository
}

func NewApplication() *Application {
	db, err := sql.Open(
		"postgres",
		helpers.GetEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
	)

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

	return &Application{
		db:             db,
		redisClient:    redisClient,
		r:              r,
		authRepository: repositories.NewDefaultAuthRepository(db, redisClient),
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

func (a *Application) RegisterControllers() {
	controller := controllers.NewAuthController(a.r, a.authRepository)
	controller.RegisterRoutes()
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
