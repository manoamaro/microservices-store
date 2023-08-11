package ports

//go:generate mockery --name InventoryService --case=snake --output ../../test/mocks
type InventoryService interface {
	Get(productId string) (uint, error)
	Add(productId string, amount uint) (uint, error)
	Subtract(productId string, amount uint) (uint, error)
	Reserve(cartId string, productId string, amount uint) (uint, error)
}
