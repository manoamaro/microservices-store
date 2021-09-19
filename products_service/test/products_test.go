package test

import (
	"testing"

	"manoamaro.github.com/products_service/internal"
	"manoamaro.github.com/products_service/internal/models"
)

func exists(arr []models.Product, f func(models.Product) bool) bool {
	for _, v := range arr {
		if f(v) {
			return true
		}
	}
	return false
}

func TestProductList(t *testing.T) {
	internal.ConnectMongoDB("mongodb://127.0.0.1:27017")

	newProduct, err := internal.DB.InsertProduct(models.Product{
		Name:        "TEST1",
		Description: "TEST TEST",
		Images:      []string{"TESTING1URL"},
		Reviews: []models.Review{
			{
				Author:  "MANOEL",
				Rating:  3,
				Comment: "SHIT",
			},
		},
		Deleted: false,
	})

	products, err := internal.DB.ListProducts()
	if err != nil {
		t.Fatal(err)
	}

	if !exists(products, func(i models.Product) bool { return i.Id == newProduct.Id }) {
		t.Error("ID of new Product not found in List")
	}

	product, err := internal.DB.FetchProduct(newProduct.Id)
	if err != nil {
		t.Fatal(err)
	}

	if product.Name != "TEST1" {
		t.Error("Name not persisted")
	}

	product.Deleted = true
	product.Name = "UPDATED"

	if updated, err := internal.DB.UpdateProduct(product.Id, *product); err != nil || !updated {
		t.Error("Not updated", err)
	}
}
