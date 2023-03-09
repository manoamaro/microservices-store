package repositories

import (
	"context"
	"database/sql"
	"github.com/manoamaro/microservices-store/order_service/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type CartRepository interface {
	Get(id uint) entities.Cart
}

type cartDBRepository struct {
	context context.Context
	db      *sql.DB
	orm     *gorm.DB
}

func NewCartDBRepository(dbUrl string) CartRepository {
	db, err := sql.Open("postgres", dbUrl)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	if err := gormDB.AutoMigrate(&entities.Cart{}, &entities.CartItem{}); err != nil {
		log.Fatal(err)
	}

	return &cartDBRepository{
		context: context.Background(),
		db:      db,
		orm:     gormDB,
	}
}

func (c *cartDBRepository) Get(id uint) entities.Cart {
	//TODO implement me
	panic("implement me")
}
