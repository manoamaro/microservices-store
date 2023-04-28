package use_cases

import (
	"fmt"
	"github.com/manoamaro/microservices-store/order_service/internal/core/domain"
	"github.com/manoamaro/microservices-store/order_service/internal/core/ports/driven"
	"github.com/manoamaro/microservices-store/order_service/internal/service"
	"strconv"
)

type GetCartUseCase interface {
	Get(id uint)
}

type GetUserCartUseCase interface {
	GetUserCart(userId string)
}

type AddItemToCartUseCase interface {
	AddItem(cartId uint, productId string, quantity uint) error
}

type UpdateItemCartUseCase interface {
	UpdateItem(cartId uint, productId string, quantity uint) error
}

type cartUseCase struct {
	repository       driven_ports.CartRepository
	productService   service.ProductService
	inventoryService service.InventoryService
}

func NewGetCartUseCase(
	repository driven_ports.CartRepository,
	productService service.ProductService,
	inventoryService service.InventoryService,
) GetCartUseCase {
	return &cartUseCase{
		repository:       repository,
		productService:   productService,
		inventoryService: inventoryService,
	}
}

func NewAddItemToCartUseCase(
	repository driven_ports.CartRepository,
	productService service.ProductService,
	inventoryService service.InventoryService,
) AddItemToCartUseCase {
	return &cartUseCase{
		repository:       repository,
		productService:   productService,
		inventoryService: inventoryService,
	}
}

func NewGetUserCartUseCase(
	repository driven_ports.CartRepository,
	productService service.ProductService,
	inventoryService service.InventoryService,
) GetUserCartUseCase {
	return &cartUseCase{
		repository:       repository,
		productService:   productService,
		inventoryService: inventoryService,
	}
}

func (c *cartUseCase) Get(id uint) {
	//TODO implement me
	panic("implement me")
}

func (c *cartUseCase) GetUserCart(userId string) {
	//TODO implement me
	panic("implement me")
}

func (c *cartUseCase) AddItem(cartId uint, productId string, quantity uint) error {
	if cart := c.repository.Get(cartId); cart.Status != domain.CartStatusOpen {
		return fmt.Errorf("cart is not open")
	} else if product, err := c.productService.Get(productId); err != nil {
		return err
	} else if _, err := c.inventoryService.Reserve(strconv.FormatInt(int64(cartId), 32), product.Id, quantity); err != nil {
		return err
	} else {
		return nil
	}
}

func (c *cartUseCase) UpdateItem(cartId uint, productId string, quantity uint) error {
	//TODO implement me
	panic("implement me")
}
