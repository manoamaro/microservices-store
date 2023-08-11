package internal

import (
	"embed"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/order_service/internal/adapters"
	"github.com/manoamaro/microservices-store/order_service/internal/application"
	ports2 "github.com/manoamaro/microservices-store/order_service/internal/ports"
	"golang.org/x/exp/slog"
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
	orderRepository  ports2.OrderRepository
	authService      infra.AuthService
	inventoryService ports2.InventoryService
	productService   ports2.ProductService
	orderApi         ports2.OrderApi
	migrator         infra.Migrator
}

func NewApplication() *Application {
	postgresUrl := helpers.GetEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/order_service?sslmode=disable")
	authUrl := helpers.GetEnv("AUTH_SERVICE_URL", "http://localhost:8080")
	inventoryUrl := helpers.GetEnv("INVENTORY_SERVICE_URL", "http://localhost:8080")
	productUrl := helpers.GetEnv("PRODUCT_SERVICE_URL", "http://localhost:8080")

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: postgresUrl,
	}), &gorm.Config{})

	if err != nil {
		slog.Error("error connecting to DB", "error", err)
	}

	engine := gin.Default()

	orderRepository, _ := adapters.NewDBOrderRepository(gormDB)
	orderService := application.NewOrderService(orderRepository)
	authService := infra.NewHttpAuthService(authUrl)
	inventoryService := adapters.NewHttpInventoryService(inventoryUrl)
	productService := adapters.NewHttpProductService(productUrl)

	return &Application{
		engine:           engine,
		migrator:         infra.NewMigrator(postgresUrl, migrationsFS),
		orderRepository:  orderRepository,
		orderApi:         adapters.NewHttpOrderApi(engine, orderService),
		authService:      authService,
		inventoryService: inventoryService,
		productService:   productService,
	}
}

func (a *Application) RegisterRoutes() {
	// Enable CORS
	a.engine.Use(cors.New(infra.CorsConfig()))
	a.engine.GET("/health", infra.HealthHandler(func() error {
		return nil
	}))
}

func (a *Application) RunMigrations() {
	if err := a.migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}

func (a *Application) Run(c chan error) {
	a.RegisterRoutes()
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
