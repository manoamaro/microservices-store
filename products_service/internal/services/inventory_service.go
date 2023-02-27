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
	infra.Service
}

func NewDefaultInventoryService(host string) InventoryService {
	return &DefaultInventoryService{
		infra.NewService(host, "InventoryService", 10, 3000),
	}
}

type Amount struct {
	Amount int `json:"amount"`
}

func (d *DefaultInventoryService) AmountOf(productId string) (int, error) {
	response, err := d.CB.Execute(func() (interface{}, error) {
		if res, err := infra.Req[Amount](d.Client, http.MethodGet, fmt.Sprintf("%s/inventory/%s", d.Host, productId), nil); err != nil {
			return nil, err
		} else {
			return res.Amount, nil
		}
	})
	return response.(int), err
}

func (d *DefaultInventoryService) Add(productId string, amount int) (int, error) {
	response, err := d.CB.Execute(func() (interface{}, error) {

		req := struct {
			ProductId string
			Amount    int
		}{productId, amount}

		if res, err := infra.Req[Amount](d.Client, http.MethodPost, fmt.Sprintf("%s/inventory/add", d.Host), req); err != nil {
			return nil, err
		} else {
			return res.Amount, nil
		}
	})
	return response.(int), err
}

func (d *DefaultInventoryService) Subtract(productId string, amount int) (int, error) {
	response, err := d.CB.Execute(func() (interface{}, error) {

		req := struct {
			ProductId string
			Amount    int
		}{productId, amount}

		if res, err := infra.Req[Amount](d.Client, http.MethodPost, fmt.Sprintf("%s/inventory/subtract", d.Host), req); err != nil {
			return nil, err
		} else {
			return res.Amount, nil
		}
	})
	return response.(int), err
}
