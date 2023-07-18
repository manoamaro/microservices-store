package service

import (
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"net/http"
)

type ProductDTO struct {
	Id          string
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       PriceDTO `json:"price"`
}

type PriceDTO struct {
	Currency string `json:"currency"`
	Price    int    `json:"price"`
}

type ProductService interface {
	Get(productId string) (ProductDTO, error)
}

type httpProductService struct {
	infra.HttpService
	getProductEndpoint *infra.Endpoint[ProductDTO]
}

func NewHttpProductService(host string) ProductService {
	service := infra.NewHttpService(host)
	return &httpProductService{
		HttpService:        *service,
		getProductEndpoint: infra.NewEndpoint[ProductDTO](service, http.MethodGet, "/public/:id", 10, 1000),
	}
}

func (h *httpProductService) Get(productId string) (ProductDTO, error) {
	res, err := h.getProductEndpoint.Start().
		WithPathParam(":id", productId).
		Execute()
	if err != nil {
		return ProductDTO{}, err
	}
	return res, nil
}
