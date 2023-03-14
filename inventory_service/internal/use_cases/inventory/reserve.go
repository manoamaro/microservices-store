package inventory

import (
	"github.com/manoamaro/microservices-store/inventory_service/internal/repository"
)

type ReserveUseCase interface {
	Reserve(cartId string, productId string, amount uint) (uint, error)
}

type reserveUseCase struct {
	repository repository.InventoryRepository
}

func NewReserveUseCase(inventoryRepository repository.InventoryRepository) ReserveUseCase {
	return &reserveUseCase{
		repository: inventoryRepository,
	}
}

func (r *reserveUseCase) Reserve(cartId string, productId string, amount uint) (uint, error) {
	return r.repository.Reserve(cartId, productId, amount)
}
