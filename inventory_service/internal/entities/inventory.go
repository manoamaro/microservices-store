package entities

import (
	"time"
)

type Inventory struct {
	ProductId string `gorm:"notNull;primarykey"`
	Amount    uint   `gorm:"notNull"`
	UpdatedAt time.Time
}
