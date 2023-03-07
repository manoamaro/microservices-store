package internal

import (
	"fmt"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"

	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
)

type Application struct {
	engine      *gin.Engine
	authService infra.AuthService
	migrator    Migrator
}

func NewApplication() *Application {
	postgresUrl := helpers.GetEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/inventory_service?sslmode=disable")
	authUrl := helpers.GetEnv("AUTH_SERVICE_URL", "http://localhost:8080")

	engine := gin.Default()
	return &Application{
		engine:      engine,
		authService: infra.NewDefaultAuthService(authUrl),
		migrator:    NewMigrator(postgresUrl),
	}
}

func (a *Application) RunMigrations() {
	if err := a.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func (a *Application) RegisterControllers() {
}

func (a *Application) Run(c chan error) {
	a.RunMigrations()
	a.RegisterControllers()

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
