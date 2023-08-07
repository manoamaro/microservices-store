package internal

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	driven2 "github.com/manoamaro/microservices-store/order_service/internal/adapters/driven"
	driver_adapters "github.com/manoamaro/microservices-store/order_service/internal/adapters/driver"
	"github.com/manoamaro/microservices-store/order_service/internal/core/application"
	driven_ports "github.com/manoamaro/microservices-store/order_service/internal/core/ports/driven"
	driver_ports "github.com/manoamaro/microservices-store/order_service/internal/core/ports/driver"
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
	engine          *gin.Engine
	orderRepository driven_ports.OrderRepository
	orderApi        driver_ports.OrderApi
	migrator        infra.Migrator
}

func NewApplication() *Application {
	postgresUrl := helpers.GetEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/order_service?sslmode=disable")
	//authUrl := helpers.GetEnv("AUTH_SERVICE_URL", "http://localhost:8080")
	//inventoryUrl := helpers.GetEnv("INVENTORY_SERVICE_URL", "http://localhost:8080")

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: postgresUrl,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()

	orderRepository, _ := driven2.NewOrderESRepository(gormDB)
	orderService := application.NewOrderService(orderRepository)

	return &Application{
		engine:          engine,
		migrator:        infra.NewMigrator(postgresUrl, migrationsFS),
		orderRepository: orderRepository,
		orderApi:        driver_adapters.NewGinOrderHandler(engine, orderService),
	}
}

func (a *Application) RegisterRoutes() {
	// Enable CORS
	a.engine.Use(cors.New(helpers.CorsConfig()))
	a.engine.Handle("GET", "/health", func(ctx *gin.Context) {

		ctx.JSON(200, gin.H{"status": "ok"})
	})
}

func (a *Application) RunMigrations() {
	if err := a.migrator.Up(); err != nil && err != migrate.ErrNoChange {
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
