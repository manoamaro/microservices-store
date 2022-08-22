package test

import (
	models2 "manoamaro.github.com/products_service/models"
	"testing"
)

func exists(arr []models2.Product, f func(models2.Product) bool) bool {
	for _, v := range arr {
		if f(v) {
			return true
		}
	}
	return false
}

func TestProductList(t *testing.T) {

}
