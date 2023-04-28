package driven_adapters

import (
	"context"
	"github.com/manoamaro/microservices-store/order_service/internal/core/domain"
	"gorm.io/gorm"
)

type orderDBRepository struct {
	context context.Context
	orm     *gorm.DB
}

func (o *orderDBRepository) Get(id uint) (domain.Order, error) {
	var result domain.Order
	tx := o.orm.Where("id = ?", id).First(&result)
	return result, tx.Error
}

func NewOrderDBRepository(gormDB *gorm.DB) driven_ports.OrderRepository {
	return &orderDBRepository{
		context: context.Background(),
		orm:     gormDB,
	}
}
