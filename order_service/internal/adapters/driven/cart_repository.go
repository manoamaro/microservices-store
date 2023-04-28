package driven_adapters

import (
	"context"
	"fmt"
	"github.com/manoamaro/microservices-store/order_service/internal/core/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type cartDBRepository struct {
	context context.Context
	orm     *gorm.DB
}

func NewCartDBRepository(gormDB *gorm.DB) driven_ports.CartRepository {
	return &cartDBRepository{
		context: context.Background(),
		orm:     gormDB,
	}
}

func (c *cartDBRepository) Get(id uint) *domain.Cart {
	var cart domain.Cart
	c.orm.
		Preload(clause.Associations).
		First(&cart, id)
	return &cart
}

func (c *cartDBRepository) CartByUserId(userId string) (*domain.Cart, error) {
	var results []domain.Cart
	tx := c.orm.
		Preload(clause.Associations).
		Where(&domain.Cart{UserId: userId, Status: domain.CartStatusOpen}).
		Find(&results)
	if tx.Error != nil {
		return nil, tx.Error
	} else if len(results) == 0 {
		return nil, fmt.Errorf("cannot find open cart for user")
	} else {
		// merge if more carts found?
		return &results[0], nil
	}
}

func (c *cartDBRepository) GetOrCreateByUserId(userId string) (*domain.Cart, error) {
	cart := domain.Cart{UserId: userId, Status: domain.CartStatusOpen}
	tx := c.orm.Debug().Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			TargetWhere: clause.Where{
				Exprs: []clause.Expression{clause.Eq{
					Column: "status",
					Value:  domain.CartStatusOpen,
				}},
			},
			DoUpdates: clause.AssignmentColumns([]string{"user_id"}),
		}, clause.Returning{},
	).Preload(clause.Associations).
		Create(&cart)

	return &cart, tx.Error
}
