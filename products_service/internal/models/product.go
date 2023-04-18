package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Prices      []Price            `bson:"prices" json:"prices"`
	Images      []string           `bson:"images" json:"images"`
	Reviews     []Review           `bson:"reviews" json:"reviews"`
	Deleted     bool               `bson:"deleted" json:"-"`
}

type Price struct {
	Currency string `bson:"currency" json:"currency"`
	Value    int    `bson:"price" json:"price"`
}

type Review struct {
	UserId  string `bson:"user_id" json:"user_id"`
	Author  string `bson:"author" json:"author"`
	Rating  int    `bson:"rating" json:"rating"`
	Comment string `bson:"comment" json:"comment"`
}

type CreateProductRequest struct {
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description" binding:"required"`
	Price       PriceRequest `json:"price" binding:"required"`
}

type PriceRequest struct {
	Currency string  `json:"currency"`
	Price    float64 `json:"price"`
}
