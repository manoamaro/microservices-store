package adapters

import (
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/order_service/internal/core/ports"
	"net/http"
)

type httpProductService struct {
	infra.HttpService
	getProductEndpoint *infra.Endpoint[ports.ProductDTO]
}

func NewHttpProductService(host string) ports.ProductService {
	service := infra.NewHttpService(host)
	return &httpProductService{
		HttpService:        *service,
		getProductEndpoint: infra.NewEndpoint[ports.ProductDTO](service, http.MethodGet, "/public/:id", 10, 1000),
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
