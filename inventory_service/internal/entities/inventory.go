package entities

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	ProductId string    `gorm:"notNull;index"`
	Operation Operation `gorm:"notNull"`
	Amount    uint      `gorm:"notNull"`
	CartId    string
}

type Operation uint8

const (
	Add Operation = iota
	Subtract
	Reserve
)
