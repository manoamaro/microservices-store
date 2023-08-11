package domain

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId          string      `gorm:"index"`
	Status          int         `gorm:"index"`
	ShippingAddress Address     `gorm:"embedded;embeddedPrefix:shipping_address_"`
	InvoiceAddress  Address     `gorm:"embedded;embeddedPrefix:invoice_address_"`
	Items           []OrderItem `gorm:"foreignKey:OrderId"`
	Total           int
}

type Address struct {
	FirstName    string
	LastName     string
	AddressLine1 string
	AddressLine2 string
	ZipCode      string
	Region       string
	State        string
	Country      string
}
