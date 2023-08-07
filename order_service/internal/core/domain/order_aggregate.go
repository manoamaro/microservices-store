package domain

import eh "github.com/looplab/eventhorizon"

type Order struct {
	UserId          string
	Status          int
	ShippingAddress Address
	InvoiceAddress  Address
	Items           []OrderItem
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

type OrderItem struct {
	ProductId string
	Quantity  uint
}

const OrderAggregateType = eh.AggregateType("order")
