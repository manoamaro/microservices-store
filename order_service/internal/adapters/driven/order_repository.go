package driven_adapters

import (
	"context"
	"github.com/manoamaro/microservices-store/commons/pkg/event_sourcing"
	"github.com/manoamaro/microservices-store/order_service/internal/core/domain"
	drivenports "github.com/manoamaro/microservices-store/order_service/internal/core/ports/driven"
	"gorm.io/gorm"
)

type orderESRepository struct {
	context context.Context
	orm     *gorm.DB
}

func NewOrderESRepository(gormDB *gorm.DB) (drivenports.OrderRepository, error) {
	if err := gormDB.AutoMigrate(&domain.OrderEvent{}); err != nil {
		return nil, err
	}

	return &orderESRepository{
		context: context.Background(),
		orm:     gormDB,
	}, nil
}

func (o *orderESRepository) Get(id string) (domain.Order, error) {
	var events []domain.OrderEvent
	result := &domain.Order{
		Entity: event_sourcing.Entity{
			ID:      id,
			Version: 0,
		},
	}

	if tx := o.orm.Where("aggregate_id = ?", id).Order("timestamp").Find(&events); tx.Error != nil {
		return *result, tx.Error
	}

	for _, event := range events {
		if event.Version != result.Version+1 {
			return *result, event_sourcing.ErrVersionConflict
		}
		result.Apply(event.Event)
		result.Version = event.Version
	}

	return *result, nil
}
