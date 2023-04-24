package repository

import (
	"context"
	"github.com/manoamaro/microservices-store/inventory_service/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const ReserveExpirationSeconds = 10

type InventoryRepository interface {
	AmountOf(productId string) (uint, error)
	Set(productId string, amount uint) (uint, error)
	Add(productId string, amount uint) (uint, error)
	Subtract(productId string, amount uint) (uint, error)
	Reserve(cartId string, productId string, amount uint) (uint, error)
}

type inventoryDBRepository struct {
	context context.Context
	ormDB   *gorm.DB
}

func NewInventoryDBRepository(gormDB *gorm.DB) InventoryRepository {
	return &inventoryDBRepository{
		context: context.Background(),
		ormDB:   gormDB,
	}
}

func (i *inventoryDBRepository) AmountOf(productId string) (uint, error) {
	var result entities.Inventory
	i.ormDB.Where("product_id = ?", productId).First(&result)
	return result.Amount, nil
}

func (i *inventoryDBRepository) Set(productId string, amount uint) (uint, error) {
	i.ormDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"amount"}),
	}).Create(&entities.Inventory{
		ProductId: productId,
		Amount:    amount,
	})
	return amount, nil
}

func (i *inventoryDBRepository) Add(productId string, a uint) (uint, error) {
	inventory := &entities.Inventory{
		ProductId: productId,
		Amount:    a,
	}
	t := i.ormDB.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "product_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"amount": gorm.Expr("inventories.amount + ?", a)}),
		},
		clause.Returning{},
	).Create(inventory)
	return inventory.Amount, t.Error
}

func (i *inventoryDBRepository) Subtract(productId string, a uint) (uint, error) {
	inventory := &entities.Inventory{
		ProductId: productId,
		Amount:    -a,
	}
	t := i.ormDB.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "product_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"amount": gorm.Expr("inventories.amount - ?", a)}),
		},
		clause.Returning{},
	).Create(inventory)
	return inventory.Amount, t.Error
}

func (i *inventoryDBRepository) Reserve(cartId string, productId string, amount uint) (uint, error) {
	return 0, nil
}
