package use_cases

import "github.com/manoamaro/microservices-store/order_service/internal/repositories"

type CartUseCase interface {
	Get(id uint)
	GetUserCart(userId string)
	AddItem(cartId uint, productId string, quantity uint)
	UpdateItem(cartId uint, productId string, quantity uint)
}

type cartUseCase struct {
	repository repositories.CartRepository
}

func (c *cartUseCase) Get(id uint) {
	//TODO implement me
	panic("implement me")
}

func (c *cartUseCase) GetUserCart(userId string) {
	//TODO implement me
	panic("implement me")
}

func (c *cartUseCase) AddItem(cartId uint, productId string, quantity uint) {
	//TODO implement me
	panic("implement me")
}

func (c *cartUseCase) UpdateItem(cartId uint, productId string, quantity uint) {
	//TODO implement me
	panic("implement me")
}
