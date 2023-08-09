package application

import (
	"github.com/manoamaro/microservices-store/commons/pkg/collections"
	"github.com/manoamaro/microservices-store/order_service/internal/core/domain"
	drivenports "github.com/manoamaro/microservices-store/order_service/internal/core/ports"
)

type OrderService struct {
	orderRepository drivenports.OrderRepository
}

func NewOrderService(orderRepository drivenports.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

func (s *OrderService) GetOrder(id uint) (*Order, error) {
	if dbOrder, err := s.orderRepository.Get(id); err != nil {
		return nil, err
	} else {
		return toDTO(dbOrder), nil
	}
}

func toDTO(dbOrder domain.Order) *Order {
	cart := Cart{
		Items: collections.MapTo(dbOrder.Items, func(i domain.OrderItem) CartItem {
			return CartItem{
				ProductId: i.ProductId,
				Quantity:  i.Quantity,
			}
		}),
		Total: 0,
	}
	return &Order{
		Id:     dbOrder.ID,
		Status: dbOrder.Status,
		Cart:   cart,
		ShippingAddress: Address{
			FirstName:    dbOrder.ShippingAddress.FirstName,
			LastName:     dbOrder.ShippingAddress.LastName,
			AddressLine1: dbOrder.ShippingAddress.AddressLine1,
			AddressLine2: dbOrder.ShippingAddress.AddressLine2,
			ZipCode:      dbOrder.ShippingAddress.ZipCode,
			Region:       dbOrder.ShippingAddress.Region,
			State:        dbOrder.ShippingAddress.State,
			Country:      dbOrder.ShippingAddress.Country,
		},
		InvoiceAddress: Address{
			FirstName:    dbOrder.InvoiceAddress.FirstName,
			LastName:     dbOrder.InvoiceAddress.LastName,
			AddressLine1: dbOrder.InvoiceAddress.AddressLine1,
			AddressLine2: dbOrder.InvoiceAddress.AddressLine2,
			ZipCode:      dbOrder.InvoiceAddress.ZipCode,
			Region:       dbOrder.InvoiceAddress.Region,
			State:        dbOrder.InvoiceAddress.State,
			Country:      dbOrder.InvoiceAddress.Country,
		},
		Total: dbOrder.Total,
	}
}
