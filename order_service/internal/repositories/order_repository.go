package repositories

import (
	"context"
	"github.com/manoamaro/microservices-store/order_service/internal/entities"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Get(id uint) entities.Cart
}

type orderDBRepository struct {
	context context.Context
	orm     *gorm.DB
}

func NewOrderDBRepository(gormDB *gorm.DB) OrderRepository {
	return &orderDBRepository{
		context: context.Background(),
		orm:     gormDB,
	}
}

func (c *orderDBRepository) Get(id uint) entities.Cart {
	//TODO implement me
	panic("implement me")
}
