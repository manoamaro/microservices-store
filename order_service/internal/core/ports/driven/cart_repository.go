package driven_ports

import "github.com/manoamaro/microservices-store/order_service/internal/core/domain"

type CartRepository interface {
	Get(id uint) *domain.Cart
	CartByUserId(userId string) (*domain.Cart, error)
	GetOrCreateByUserId(userId string) (*domain.Cart, error)
}
