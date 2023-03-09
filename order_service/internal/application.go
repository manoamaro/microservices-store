package internal

import (
	"fmt"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/order_service/internal/repositories"
	"github.com/manoamaro/microservices-store/order_service/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
)

type Application struct {
	engine           *gin.Engine
	authService      infra.AuthService
	inventoryService service.InventoryService
	cartRepository   repositories.CartRepository
	controllers      []infra.Controller
}

func NewApplication() *Application {
	postgresUrl := helpers.GetEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/order_service?sslmode=disable")
	authUrl := helpers.GetEnv("AUTH_SERVICE_URL", "http://localhost:8080")
	inventoryUrl := helpers.GetEnv("INVENTORY_SERVICE_URL", "http://localhost:8080")

	engine := gin.Default()
	return &Application{
		engine:           engine,
		authService:      infra.NewDefaultAuthService(authUrl),
		inventoryService: service.NewHttpInventoryService(inventoryUrl),
		cartRepository:   repositories.NewCartDBRepository(postgresUrl),
	}
}

func (a *Application) RegisterControllers() {
	for _, controller := range a.controllers {
		controller.RegisterRoutes()
	}
}

func (a *Application) Run(c chan error) {
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
