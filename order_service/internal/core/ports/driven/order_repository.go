package driven_ports

import "github.com/manoamaro/microservices-store/order_service/internal/core/domain"

type OrderRepository interface {
	Get(id uint) (domain.Order, error)
}
