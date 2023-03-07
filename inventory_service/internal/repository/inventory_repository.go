package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/manoamaro/microservices-store/inventory_service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	AmountOf(productId string) (amount uint)
	Add(productId string, a uint) (amount uint)
	Subtract(productId string, a uint) (amount uint)
}

type DefaultInventoryRepository struct {
	context context.Context
	db      *sql.DB
	ormDB   *gorm.DB
}

func NewDefaultInventoryRepository(dbUrl string) InventoryRepository {
	db, err := sql.Open("postgres", dbUrl)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return &DefaultInventoryRepository{
		context: context.Background(),
		db:      db,
		ormDB:   gormDB,
	}
}

func (i *DefaultInventoryRepository) AmountOf(productId string) (amount uint) {
	return amountOf(i.ormDB, productId)
}

func amountOf(tx *gorm.DB, productId string) (amount uint) {
	var transactions []models.Transaction
	tx.Where("product_id = ?", productId).Order("created_at asc").Find(&transactions)
	amount = 0
	for _, t := range transactions {
		switch t.Operation {
		case models.Add:
			amount += t.Amount
		case models.Subtract:
			amount -= t.Amount
		}
	}

	return amount
}

func (i *DefaultInventoryRepository) Add(productId string, a uint) (amount uint) {
	newTransaction := models.Transaction{
		ProductId: productId,
		Amount:    a,
		Operation: models.Add,
	}

	amount = 0

	i.ormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newTransaction).Error; err != nil {
			return err
		}
		amount = amountOf(tx, productId)
		return nil
	})

	return amount
}

func (i *DefaultInventoryRepository) Subtract(productId string, a uint) (amount uint) {
	newTransaction := models.Transaction{
		ProductId: productId,
		Amount:    a,
		Operation: models.Subtract,
	}

	amount = 0

	i.ormDB.Transaction(func(tx *gorm.DB) error {
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

	return amount
}
