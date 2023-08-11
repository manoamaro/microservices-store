package ports

import (
	"github.com/manoamaro/microservices-store/order_service/internal/domain"
)

//go:generate mockery --name OrderRepository --case=snake --output ../../test/mocks
type OrderRepository interface {
	Get(id uint) (domain.Order, error)
}
