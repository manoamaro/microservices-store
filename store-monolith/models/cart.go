package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	Id        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserId    primitive.ObjectID   `bson:"user_id" json:"user_id"`
	Products  []primitive.ObjectID `bson:"product_ids" json:"product_ids"`
	Discounts []int                `bson:"discounts" json:"discounts"`
	Currency  string               `bson:"currency" json:"currency"`
	Total     int                  `bson:"total" json:"total"`
}
