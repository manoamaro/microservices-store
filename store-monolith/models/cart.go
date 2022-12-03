package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId   primitive.ObjectID `bson:"user_id" json:"user_id"`
	Products []struct {
		Id       primitive.ObjectID `bson:"id,omitempty" json:"id"`
		Value    int                `bson:"value" json:"value"`
		Quantity int                `bson:"quantity" json:"quantity"`
	} `bson:"products" json:"products"`
	Discounts []int   `bson:"discounts" json:"discounts"`
	Currency  string  `bson:"currency" json:"currency"`
	Total     float32 `bson:"total" json:"total"`
}
