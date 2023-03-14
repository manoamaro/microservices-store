package repositories

import (
	"context"
	"fmt"
	"github.com/manoamaro/microservices-store/order_service/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartRepository interface {
	Get(id uint) *entities.Cart
	CartByUserId(userId string) (*entities.Cart, error)
	GetOrCreateByUserId(userId string) (*entities.Cart, error)
}

type cartDBRepository struct {
	context context.Context
	orm     *gorm.DB
}

func NewCartDBRepository(gormDB *gorm.DB) CartRepository {
	return &cartDBRepository{
		context: context.Background(),
		orm:     gormDB,
	}
}

func (c *cartDBRepository) Get(id uint) *entities.Cart {
	var cart entities.Cart
	c.orm.
		Preload(clause.Associations).
		First(&cart, id)
	return &cart
}

func (c *cartDBRepository) CartByUserId(userId string) (*entities.Cart, error) {
	var results []entities.Cart
	tx := c.orm.
		Preload(clause.Associations).
		Where(&entities.Cart{UserId: userId, Status: entities.CartStatusOpen}).
		Find(&results)
	if tx.Error != nil {
		return nil, tx.Error
	} else if len(results) == 0 {
		return nil, fmt.Errorf("cannot find open cart for user")
	} else {
		// merge if more carts found?
		return &results[0], nil
	}
}

func (c *cartDBRepository) GetOrCreateByUserId(userId string) (*entities.Cart, error) {
	var cart entities.Cart
	tx := c.orm.
		Preload(clause.Associations).
		Where(&entities.Cart{UserId: userId, Status: entities.CartStatusOpen}).
		FirstOrCreate(&cart)
	return &cart, tx.Error
}
