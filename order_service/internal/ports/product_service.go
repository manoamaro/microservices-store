package ports

//go:generate mockery --name ProductService --case=snake --output ../../test/mocks
type ProductService interface {
	Get(productId string) (ProductDTO, error)
}

type ProductDTO struct {
	Id          string
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       PriceDTO `json:"price"`
}

type PriceDTO struct {
	Currency string `json:"currency"`
	Price    int    `json:"price"`
}
