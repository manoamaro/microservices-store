package adapters

import (
	"context"
	"github.com/manoamaro/microservices-store/order_service/internal/domain"
	drivenports "github.com/manoamaro/microservices-store/order_service/internal/ports"
	"gorm.io/gorm"
)

type dbOrderRepository struct {
	context context.Context
	orm     *gorm.DB
}

func NewDBOrderRepository(gormDB *gorm.DB) (drivenports.OrderRepository, error) {
	//if err := gormDB.AutoMigrate(&domain.Order{}, &domain.OrderItem{}); err != nil {
	//	return nil, err
	//}

	return &dbOrderRepository{
		context: context.Background(),
		orm:     gormDB,
	}, nil
}

func (o *dbOrderRepository) Get(id uint) (domain.Order, error) {
	var result domain.Order
	if err := o.orm.First(&result, id).Error; err != nil {
		return domain.Order{}, err
	}
	return result, nil
}
