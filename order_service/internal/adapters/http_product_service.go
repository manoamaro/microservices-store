package adapters

import (
	"github.com/manoamaro/microservices-store/commons/pkg/http_client"
	"github.com/manoamaro/microservices-store/order_service/internal/ports"
	"net/http"
)

type httpProductService struct {
	http_client.HttpClient
	getProductEndpoint *http_client.Endpoint[any, ports.ProductDTO]
}

func NewHttpProductService(host string) ports.ProductService {
	service := http_client.NewHttpClient(host)
	return &httpProductService{
		HttpClient:         *service,
		getProductEndpoint: http_client.NewEndpoint[any, ports.ProductDTO](service, http.MethodGet, "/public/:id", 10, 1000),
	}
}

func (h *httpProductService) Get(productId string) (ports.ProductDTO, error) {
	res, err := h.getProductEndpoint.Start().
		WithPathParam(":id", productId).
		Execute()
	if err != nil {
		return ports.ProductDTO{}, err
	}
	return res, nil
}
