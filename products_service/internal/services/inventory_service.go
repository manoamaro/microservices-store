package services

import (
	"fmt"
	"net/http"

	"github.com/manoamaro/microservices-store/commons/pkg/infra"
)

type InventoryService interface {
	infra.IService
	AmountOf(productId string) (int, error)
	Add(productId string, amount int) (int, error)
	Subtract(productId string, amount int) (int, error)
}

type DefaultInventoryService struct {
	*infra.Service
	amountOfEndpoint *infra.Endpoint[int]
	addEndpoint      *infra.Endpoint[int]
	subtractEndpoint *infra.Endpoint[int]
}

func NewDefaultInventoryService(host string) InventoryService {
	service := infra.NewService(host)
	return &DefaultInventoryService{
		Service:          service,
		amountOfEndpoint: infra.NewEndpoint[int](service, http.MethodGet, 10, 3000),
		addEndpoint:      infra.NewEndpoint[int](service, http.MethodPost, 10, 3000),
		subtractEndpoint: infra.NewEndpoint[int](service, http.MethodPost, 10, 3000),
	}
}

type Amount struct {
	Amount int `json:"amount"`
}

func (d *DefaultInventoryService) AmountOf(productId string) (int, error) {
	return d.amountOfEndpoint.Execute(fmt.Sprintf("/inventory/%s", productId), nil, nil)
}

func (d *DefaultInventoryService) Add(productId string, amount int) (int, error) {
	req := struct {
		ProductId string
		Amount    int
	}{productId, amount}

	return d.amountOfEndpoint.Execute("/inventory/add", nil, req)
}

func (d *DefaultInventoryService) Subtract(productId string, amount int) (int, error) {
	req := struct {
		ProductId string
		Amount    int
	}{productId, amount}
	return d.amountOfEndpoint.Execute("/inventory/subtract", nil, req)
}
