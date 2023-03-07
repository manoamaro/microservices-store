package inventory

import (
	"github.com/manoamaro/microservices-store/inventory_service/internal/repository"
)

type GetUseCase interface {
	Get(productId string) (uint, error)
}

type getUseCase struct {
	repository repository.InventoryRepository
}

func NewGetUseCase(inventoryRepository repository.InventoryRepository) GetUseCase {
	return &getUseCase{
		repository: inventoryRepository,
	}
}

func (r *getUseCase) Get(productId string) (uint, error) {
	amount := r.repository.AmountOf(productId)
	return amount, nil
}
