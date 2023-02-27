package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"

	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/inventory_service/internal/controller"
	"github.com/manoamaro/microservices-store/inventory_service/internal/repository"
)

type Application struct {
	db                  *sql.DB
	r                   *gin.Engine
	inventoryRepository repository.InventoryRepository
}

func NewApplication() *Application {
	db, err := sql.Open(
		"postgres",
		helpers.GetEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/inventory?sslmode=disable"),
	)

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	return &Application{
		db:                  db,
		r:                   r,
		inventoryRepository: repository.NewDefaultInventoryRepository(db),
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
	controller := controller.NewInventoryController(a.r, a.inventoryRepository)
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
