package entities

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserId string     `gorm:"notNull;index"`
	Status string     `gorm:"notNull;index"`
	Items  []CartItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Total  uint       `gorm:"notNull;check:total >= 0"`
}

type CartItem struct {
	gorm.Model
	CartID    uint   `gorm:"notNull;index"`
	ProductId string `gorm:"notNull;index"`
	Quantity  uint   `gorm:"notNull;check:quantity > 0"`
}
