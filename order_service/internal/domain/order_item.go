package domain

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderId   uint   `gorm:"index"`
	Order     Order  `gorm:"foreignKey:OrderId"`
	ProductId string `gorm:"index"`
	Quantity  uint
}
