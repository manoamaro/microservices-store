package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	ProductId string    `gorm:"notNull;index"`
	Operation Operation `gorm:"notNull"`
	Amount    uint      `gorm:"notNull"`
}

type Operation uint8

const (
	Add Operation = iota
	Subtract
	Reserve
)
