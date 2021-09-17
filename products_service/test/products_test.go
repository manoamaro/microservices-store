package test

import (
	"testing"

	"manoamaro.github.com/products_service/internal"
	"manoamaro.github.com/products_service/internal/models"
)

func TestProductList(t *testing.T) {
	internal.ConnectMongoDB("mongodb://127.0.0.1:27017")

	internal.DB.InsertProduct(models.Product{
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

	if products[0].Name != "TEST1" {
		t.Error("Name not persisted")
	}

	product, err := internal.DB.FetchProduct(products[0].Id)
	if err != nil {
		t.Fatal(err)
	}

	if product.Name != "TEST1" {
		t.Error("Name not persisted")
	}
}
