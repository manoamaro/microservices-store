package application

type Cart struct {
	Items []CartItem `json:"items"`
	Total uint       `json:"total"`
}

type CartItem struct {
	ProductId string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
}
