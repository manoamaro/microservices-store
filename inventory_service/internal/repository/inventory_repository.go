package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/manoamaro/microservices-store/inventory_service/internal/entities"
	"gorm.io/gorm"
	"time"
)

type InventoryRepository interface {
	AmountOf(productId string) (amount uint)
	Add(productId string, a uint) (amount uint)
	Subtract(productId string, a uint) (amount uint)
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
			expiresAt := t.CreatedAt.Add(time.Hour * 1)
			if t.CreatedAt.Before(expiresAt) {
				amount -= t.Amount
			}
		}
	}

	return amount
}

func (i *inventoryDBRepository) Add(productId string, a uint) (amount uint) {
	newTransaction := entities.Transaction{
		ProductId: productId,
		Amount:    a,
		Operation: entities.Add,
	}

	amount = 0

	err := i.ormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newTransaction).Error; err != nil {
			return err
		}
		amount = amountOf(tx, productId)
		return nil
	})

	if err != nil {
		return 0
	}

	return amount
}

func (i *inventoryDBRepository) Subtract(productId string, a uint) (amount uint) {
	newTransaction := entities.Transaction{
		ProductId: productId,
		Amount:    a,
		Operation: entities.Subtract,
	}

	amount = 0

	err := i.ormDB.Transaction(func(tx *gorm.DB) error {
		amount = amountOf(tx, productId)

		if amount < a {
			return errors.New("not enough inventory")
		}

		if err := tx.Create(&newTransaction).Error; err != nil {
			return err
		}

		amount = amountOf(tx, productId)

		return nil
	})

	if err != nil {
		return 0
	}

	return amount
}

func (i *inventoryDBRepository) Reserve(cartId string, productId string, amount uint) (uint, error) {
	tx := i.ormDB.Exec(
		"WITH amount_of AS( SELECT SUM( CASE WHEN operation = 0 THEN amount WHEN operation = 1 THEN -amount WHEN operation = 2 AND created_at < NOW() + interval '10' minute THEN -amount END) AS total FROM transactions WHERE product_id = @productId ) INSERT INTO transactions (product_id, operation, amount, cart_id) SELECT @productId, 2, @amount, @cartId FROM amount_of WHERE total >= @amount",
		sql.Named("productId", productId),
		sql.Named("cartId", cartId),
		sql.Named("amount", amount),
	)
	if tx.Error != nil {
		return 0, tx.Error
	} else if tx.RowsAffected == 0 {
		return 0, fmt.Errorf("cannot reserve")
	} else {
		return i.AmountOf(productId), nil
	}
}
