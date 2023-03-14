package controller

import (
	"fmt"
	"github.com/manoamaro/microservices-store/products_service/internal/models"
	"github.com/samber/lo"
)

type ProductDTO struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Price       ProductPriceDTO    `json:"price"`
	Reviews     []ProductReviewDTO `json:"reviews"`
}

func FromProduct(product models.Product, currency string) (ProductDTO, error) {
	price, found := lo.Find(product.Prices, func(item models.Price) bool {
		return item.Currency == currency
	})
	if !found {
		return ProductDTO{}, fmt.Errorf("price with currency not found")
	}

	reviews := lo.Map[models.Review, ProductReviewDTO](product.Reviews, func(item models.Review, index int) ProductReviewDTO {
		return ProductReviewDTO{
			UserId:  item.UserId,
			Rating:  item.Rating,
			Comment: item.Comment,
		}
	})

	return ProductDTO{
		Id:          product.Id.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Price: ProductPriceDTO{
			Currency: price.Currency,
			Value:    price.Value,
		},
		Reviews: reviews,
	}, nil
}

type ProductPriceDTO struct {
	Currency string `json:"currency"`
	Value    int    `json:"value"`
}

type ProductReviewDTO struct {
	UserId  string `json:"user_id"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}
