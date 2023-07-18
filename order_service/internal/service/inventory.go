package service

import (
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"net/http"
)

type InventoryService interface {
	Get(productId string) (uint, error)
	Add(productId string, amount uint) (uint, error)
	Subtract(productId string, amount uint) (uint, error)
	Reserve(cartId string, productId string, amount uint) (uint, error)
}

type httpInventoryService struct {
	*infra.HttpService
	reserveEndpoint *infra.Endpoint[uint]
}

func NewHttpInventoryService(host string) InventoryService {
	service := infra.NewHttpService(host)
	return &httpInventoryService{
		HttpService:     service,
		reserveEndpoint: infra.NewEndpoint[uint](service, http.MethodPost, "/public/reserve", 10, 1000),
	}
}

func (h *httpInventoryService) Get(productId string) (uint, error) {
	//TODO implement me
	panic("implement me")
}

func (h *httpInventoryService) Add(productId string, amount uint) (uint, error) {
	//TODO implement me
	panic("implement me")
}

func (h *httpInventoryService) Subtract(productId string, amount uint) (uint, error) {
	//TODO implement me
	panic("implement me")
}

func (h *httpInventoryService) Reserve(cartId string, productId string, amount uint) (uint, error) {
	req := struct {
		ProductId string `json:"product_id"`
		CartId    string `json:"cart_id"`
		Amount    uint   `json:"amount"`
	}{
		ProductId: productId,
		CartId:    cartId,
		Amount:    amount,
	}

	res, err := h.reserveEndpoint.Start().
		WithBody(req).
		Execute()

	if err != nil {
		return 0, err
	}

	return res, nil
}
