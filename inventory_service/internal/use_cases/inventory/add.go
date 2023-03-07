package inventory

import (
	"github.com/manoamaro/microservices-store/inventory_service/internal/repository"
)

type AddUseCase interface {
	Add(productId string, amount uint) (uint, error)
}

type addUseCase struct {
	repository repository.InventoryRepository
}

func NewAddUseCase(inventoryRepository repository.InventoryRepository) AddUseCase {
	return &addUseCase{
		repository: inventoryRepository,
	}
}

func (r *addUseCase) Add(productId string, amount uint) (uint, error) {
	amountAfter := r.repository.Add(productId, amount)
	return amountAfter, nil
}
