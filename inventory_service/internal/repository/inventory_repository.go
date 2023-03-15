package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/manoamaro/microservices-store/inventory_service/internal/entities"
	"gorm.io/gorm"
	"time"
)

const ReserveExpirationSeconds = 10

type InventoryRepository interface {
	AmountOf(productId string) (amount uint)
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

func (i *inventoryDBRepository) AmountOf(productId string) (amount uint) {
	return amountOf(i.ormDB, productId)
}

func amountOf(tx *gorm.DB, productId string) (amount uint) {
	var transactions []entities.Transaction
	tx.Where("product_id = ?", productId).Order("created_at asc").Find(&transactions)

	amount = 0

	for _, t := range transactions {
		switch t.Operation {
		case entities.Add:
			amount += t.Amount
		case entities.Subtract:
			amount -= t.Amount
		case entities.Reserve:
			deadline := t.CreatedAt.Add(time.Second * ReserveExpirationSeconds)
			if time.Now().Before(deadline) {
				amount -= t.Amount
			}
		}
	}

	return amount
}

func (i *inventoryDBRepository) Add(productId string, a uint) (uint, error) {
	newTransaction := entities.Transaction{
		ProductId: productId,
		Amount:    a,
		Operation: entities.Add,
	}

	tx := i.ormDB.Create(&newTransaction)
	if tx.Error != nil {
		return 0, tx.Error
	}

	amount := amountOf(i.ormDB, productId)

	return amount, nil
}

func (i *inventoryDBRepository) subtract(productId string, cartId string, amount uint, operation entities.Operation) (uint, error) {
	tx := i.ormDB.Exec(
		`WITH amount_of AS(
				SELECT SUM(
					CASE
						WHEN operation = 0 THEN amount
						WHEN operation = 1 THEN -amount
						WHEN operation = 2 AND created_at > @expires_at THEN -amount
					END
				) AS total FROM transactions WHERE product_id = @productId
			) INSERT INTO transactions (product_id, operation, amount, cart_id)
				SELECT @productId, @op, @amount, @cartId
				FROM amount_of
				WHERE total >= @amount`,
		sql.Named("expires_at", time.Now().Add(-(time.Second*ReserveExpirationSeconds))),
		sql.Named("productId", productId),
		sql.Named("cartId", cartId),
		sql.Named("amount", amount),
		sql.Named("op", int(operation)),
	)
	if tx.Error != nil {
		return 0, tx.Error
	} else if tx.RowsAffected == 0 {
		return 0, fmt.Errorf("cannot subtract")
	} else {
		return i.AmountOf(productId), nil
	}
}

func (i *inventoryDBRepository) Subtract(productId string, a uint) (uint, error) {
	return i.subtract(productId, "", a, entities.Subtract)
}

func (i *inventoryDBRepository) Reserve(cartId string, productId string, amount uint) (uint, error) {
	return i.subtract(productId, cartId, amount, entities.Reserve)
}
