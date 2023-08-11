package adapters

import (
	"github.com/manoamaro/microservices-store/commons/pkg/http_client"
	"github.com/manoamaro/microservices-store/order_service/internal/ports"
	"net/http"
)

type httpInventoryService struct {
	*http_client.HttpClient
	reserveEndpoint *http_client.Endpoint[InventoryVerifyRequest, uint]
}

func NewHttpInventoryService(host string) ports.InventoryService {
	service := http_client.NewHttpClient(host)
	return &httpInventoryService{
		HttpClient:      service,
		reserveEndpoint: http_client.NewEndpoint[InventoryVerifyRequest, uint](service, http.MethodPost, "/public/reserve", 10, 1000),
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
	req := InventoryVerifyRequest{
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

type InventoryVerifyRequest struct {
	ProductId string `json:"product_id"`
	CartId    string `json:"cart_id"`
	Amount    uint   `json:"amount"`
}
