package inventory

import (
	"github.com/manoamaro/microservices-store/inventory_service/internal/repository"
)

type SubtractUseCase interface {
	Subtract(productId string, amount uint) (uint, error)
}

type subtractUseCase struct {
	repository repository.InventoryRepository
}

func NewSubtractUseCase(inventoryRepository repository.InventoryRepository) SubtractUseCase {
	return &subtractUseCase{
		repository: inventoryRepository,
	}
}

func (r *subtractUseCase) Subtract(productId string, amount uint) (uint, error) {
	return r.repository.Subtract(productId, amount)
}
