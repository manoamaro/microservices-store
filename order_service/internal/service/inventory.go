package service

import "github.com/manoamaro/microservices-store/commons/pkg/infra"

type InventoryService interface {
	infra.IService
	Get(productId string) (uint, error)
	Add(productId string, amount uint) (uint, error)
	Subtract(productId string, amount uint) (uint, error)
}

type httpInventoryService struct {
	infra.Service
}

func NewHttpInventoryService(host string) InventoryService {
	return &httpInventoryService{
		Service: *infra.NewService(host),
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
