package services

import (
	"net/http"

	"github.com/manoamaro/microservices-store/commons/pkg/infra"
)

//go:generate mockery --name InventoryService --case=snake --output ../../test/mocks
type InventoryService interface {
	AmountOf(productId string) (int, error)
	Add(productId string, amount int) (int, error)
	Subtract(productId string, amount int) (int, error)
}

type DefaultInventoryService struct {
	*infra.HttpService
	amountOfEndpoint *infra.Endpoint[int]
	addEndpoint      *infra.Endpoint[int]
	subtractEndpoint *infra.Endpoint[int]
}

func NewDefaultInventoryService(host string) InventoryService {
	service := infra.NewHttpService(host)
	return &DefaultInventoryService{
		HttpService:      service,
		amountOfEndpoint: infra.NewEndpoint[int](service, http.MethodGet, "/inventory/:productId", 10, 3000),
		addEndpoint:      infra.NewEndpoint[int](service, http.MethodPost, "/inventory/add", 10, 3000),
		subtractEndpoint: infra.NewEndpoint[int](service, http.MethodPost, "/inventory/subtract", 10, 3000),
	}
}

type Amount struct {
	Amount int `json:"amount"`
}

func (d *DefaultInventoryService) AmountOf(productId string) (int, error) {
	return d.amountOfEndpoint.Start().
		WithPathParam(":productId", productId).
		Execute()
}

func (d *DefaultInventoryService) Add(productId string, amount int) (int, error) {
	req := struct {
		ProductId string
		Amount    int
	}{productId, amount}
	return d.addEndpoint.Start().WithBody(req).Execute()
}

func (d *DefaultInventoryService) Subtract(productId string, amount int) (int, error) {
	req := struct {
		ProductId string
		Amount    int
	}{productId, amount}
	return d.subtractEndpoint.Start().WithBody(req).Execute()
}
