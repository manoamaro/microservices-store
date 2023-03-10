package entities

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId          string  `gorm:"notNull;index"`
	Status          int     `gorm:"notNull;index"`
	CartId          uint    `gorm:"notNull;index"`
	Cart            Cart    `gorm:"notNull"`
	ShippingAddress Address `gorm:"embedded;embeddedPrefix:shipping_address_"`
	InvoiceAddress  Address `gorm:"embedded;embeddedPrefix:invoice_address_"`
	Total           int
}

type Address struct {
	FirstName    string `gorm:"notNull"`
	LastName     string `gorm:"notNull"`
	AddressLine1 string `gorm:"notNull"`
	AddressLine2 string
	ZipCode      string `gorm:"notNull"`
	Region       string
	State        string `gorm:"notNull"`
	Country      string `gorm:"notNull"`
}
