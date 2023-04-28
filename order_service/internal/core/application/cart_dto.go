package application

type Cart struct {
	Id     uint       `json:"id"`
	Status int        `json:"status"`
	Items  []CartItem `json:"items"`
	Total  uint       `json:"total"`
}

type CartItem struct {
	ProductId string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
}
