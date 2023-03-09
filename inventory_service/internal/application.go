package internal

import (
	"embed"
	"fmt"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/inventory_service/internal/use_cases/inventory"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"

	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/inventory_service/internal/controller"
	"github.com/manoamaro/microservices-store/inventory_service/internal/repository"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Application struct {
	engine              *gin.Engine
	inventoryRepository repository.InventoryRepository
	authService         infra.AuthService
	migrator            infra.Migrator
	controllers         []infra.Controller
}

func NewApplication() *Application {
	postgresUrl := helpers.GetEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/inventory_service?sslmode=disable")
	authUrl := helpers.GetEnv("AUTH_SERVICE_URL", "http://localhost:8080")

	engine := gin.Default()
	inventoryRepository := repository.NewInventoryDBRepository(postgresUrl)
	authService := infra.NewDefaultAuthService(authUrl)
	return &Application{
		engine:              engine,
		inventoryRepository: inventoryRepository,
		authService:         authService,
		migrator:            infra.NewMigrator(postgresUrl, migrationsFS),
		controllers: []infra.Controller{
			controller.NewInventoryController(
				engine,
				authService,
				inventory.NewGetUseCase(inventoryRepository),
				inventory.NewAddUseCase(inventoryRepository),
				inventory.NewSubtractUseCase(inventoryRepository),
			),
		},
	}
}

func (a *Application) RunMigrations() {
	if err := a.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func (a *Application) RegisterControllers() {
	for _, _controller := range a.controllers {
		_controller.RegisterRoutes()
	}
}

func (a *Application) Run(c chan error) {
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
