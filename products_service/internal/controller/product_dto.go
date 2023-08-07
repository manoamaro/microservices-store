package controller

import (
	"fmt"
	"github.com/manoamaro/microservices-store/commons/pkg/collections"
	"github.com/manoamaro/microservices-store/products_service/internal/models"
	"github.com/samber/lo"
)

type ProductDTO struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Images      []ProductImageDTO  `json:"images"`
	Price       ProductPriceDTO    `json:"price"`
	Reviews     []ProductReviewDTO `json:"reviews"`
}

type ProductImageDTO struct {
	Id          string `json:"id"`
	Url         string `json:"url"`
	Description string `json:"description"`
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

type ProductAdminDTO struct {
	ProductDTO
	Prices  []ProductPriceDTO `json:"prices"`
	Deleted bool              `json:"deleted"`
}

func FromProduct(product models.Product, currency string, host string) (ProductDTO, error) {
	price, found := lo.Find(product.Prices, func(item models.Price) bool {
		return item.Currency == currency
	})
	if !found && len(currency) > 0 {
		return ProductDTO{}, fmt.Errorf("price with currency not found")
	}

	reviews := lo.Map[models.Review, ProductReviewDTO](product.Reviews, func(item models.Review, index int) ProductReviewDTO {
		return ProductReviewDTO{
			UserId:  item.UserId,
			Rating:  item.Rating,
			Comment: item.Comment,
		}
	})

	images := collections.MapTo(product.Images, func(i string) ProductImageDTO {
		return ProductImageDTO{
			Id:          i,
			Url:         fmt.Sprintf("%s/public/assets/%s", host, i),
			Description: i,
		}
	})

	return ProductDTO{
		Id:          product.Id.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Images:      images,
		Price: ProductPriceDTO{
			Currency: price.Currency,
			Value:    price.Value,
		},
		Reviews: reviews,
	}, nil
}

func FromProductAdmin(product models.Product, host string) ProductAdminDTO {
	productDTO, _ := FromProduct(product, "", host)
	prices := lo.Map[models.Price, ProductPriceDTO](product.Prices, func(item models.Price, index int) ProductPriceDTO {
		return ProductPriceDTO{
			Currency: item.Currency,
			Value:    item.Value,
		}
	})
	return ProductAdminDTO{
		productDTO,
		prices,
		product.Deleted,
	}
}
