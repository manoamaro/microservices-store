package repositories

import (
	"context"
	"database/sql"
	"github.com/manoamaro/microservices-store/order_service/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type OrderRepository interface {
	Get(id uint) entities.Cart
}

type orderDBRepository struct {
	context context.Context
	db      *sql.DB
	orm     *gorm.DB
}

func NewOrderDBRepository(dbUrl string) OrderRepository {
	db, err := sql.Open("postgres", dbUrl)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	if err := gormDB.AutoMigrate(&entities.Order{}); err != nil {
		log.Fatal(err)
	}

	return &orderDBRepository{
		context: context.Background(),
		db:      db,
		orm:     gormDB,
	}
}

func (c *orderDBRepository) Get(id uint) entities.Cart {
	//TODO implement me
	panic("implement me")
}
