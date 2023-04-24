package internal

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/order_service/internal/controllers"
	"github.com/manoamaro/microservices-store/order_service/internal/repositories"
	"github.com/manoamaro/microservices-store/order_service/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Application struct {
	engine           *gin.Engine
	authService      infra.AuthService
	inventoryService service.InventoryService
	cartRepository   repositories.CartRepository
	orderRepository  repositories.OrderRepository
	controllers      []infra.Controller
	migrator         infra.Migrator
}

func NewApplication() *Application {
	postgresUrl := helpers.GetEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/order_service?sslmode=disable")
	authUrl := helpers.GetEnv("AUTH_SERVICE_URL", "http://localhost:8080")
	inventoryUrl := helpers.GetEnv("INVENTORY_SERVICE_URL", "http://localhost:8080")

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: postgresUrl,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()

	authService := infra.NewDefaultAuthService(authUrl)
	inventoryService := service.NewHttpInventoryService(inventoryUrl)

	orderRepository := repositories.NewOrderDBRepository(gormDB)
	cartRepository := repositories.NewCartDBRepository(gormDB)

	return &Application{
		engine:           engine,
		migrator:         infra.NewMigrator(postgresUrl, migrationsFS),
		authService:      authService,
		inventoryService: inventoryService,
		cartRepository:   cartRepository,
		orderRepository:  orderRepository,
		controllers: []infra.Controller{
			controllers.NewOrderController(engine, authService, orderRepository),
		},
	}
}

func (a *Application) RegisterControllers() {
	// Enable CORS
	a.engine.Use(cors.New(helpers.CorsConfig()))
	for _, controller := range a.controllers {
		controller.RegisterRoutes()
	}
}

func (a *Application) RunMigrations() {
	if err := a.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func (a *Application) Run(c chan error) {
	a.RegisterControllers()
	a.RunMigrations()

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
