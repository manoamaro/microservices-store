package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/manoamaro/microservices-store/order_service/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type CartRepository interface {
	Get(id uint) *entities.Cart
	GetOpenByUserId(userId string) (*entities.Cart, error)
	GetOpenOrCreateByUserId(userId string) (*entities.Cart, error)
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

func (c *cartDBRepository) Get(id uint) *entities.Cart {
	//TODO implement me
	panic("implement me")
}

func (c *cartDBRepository) GetOpenByUserId(userId string) (*entities.Cart, error) {
	var results []entities.Cart
	tx := c.orm.Preload(clause.Associations).
		Where("user_id = ? AND status = ?", userId, entities.CartStatusOpen).
		Find(&results)
	if tx.Error != nil {
		return nil, tx.Error
	} else if len(results) == 0 {
		return nil, fmt.Errorf("cannot find open cart for user")
	} else {
		// merge if more carts found
		return &results[0], nil
	}
}

func (c *cartDBRepository) GetOpenOrCreateByUserId(userId string) (*entities.Cart, error) {
	var cart entities.Cart
	sql := c.orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tx.Raw("")
		return tx.Model(&entities.Cart{}).Select("1").Where(&entities.Cart{UserId: userId, Status: entities.CartStatusOpen})
	})
	fmt.Println(sql)
	return &cart, nil
}
