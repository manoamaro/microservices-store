package application

type Order struct {
	Id              uint    `json:"id"`
	Status          int     `json:"status"`
	Cart            Cart    `json:"cart"`
	ShippingAddress Address `json:"shipping_address"`
	InvoiceAddress  Address `json:"invoice_address"`
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
