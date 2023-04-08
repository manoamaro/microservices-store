package internal

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/manoamaro/microservices-store/auth_service/internal/controllers"
	"github.com/manoamaro/microservices-store/auth_service/internal/use_cases"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/golang-migrate/migrate/v4"
	"github.com/manoamaro/microservices-store/auth_service/internal/repositories"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Application struct {
	db             *sql.DB
	redisClient    *redis.Client
	engine         *gin.Engine
	migrator       infra.Migrator
	authRepository repositories.AuthRepository
	controllers    []infra.Controller
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

	engine := gin.Default()
	engine.Use(cors.Default())

	authRepository := repositories.NewDefaultAuthRepository(db, redisClient)
	authController := controllers.NewAuthController(
		engine,
		use_cases.NewSignInUseCase(authRepository),
		use_cases.NewSignUpUseCase(authRepository),
		use_cases.NewVerifyUseCase(authRepository),
		use_cases.NewRefreshTokenUseCase(authRepository),
	)

	return &Application{
		db:             db,
		redisClient:    redisClient,
		engine:         engine,
		migrator:       infra.NewMigrator(postgresUrl, migrationsFS),
		authRepository: authRepository,
		controllers:    []infra.Controller{authController},
	}
}

func (a *Application) RunMigrations() {
	if err := a.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func (a *Application) Run(c chan error) {
	for _, _controller := range a.controllers {
		_controller.RegisterRoutes()
	}

	port := helpers.GetEnv("PORT", "8080")

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			c <- err
		}
	}()
}
