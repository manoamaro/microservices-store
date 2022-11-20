package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	Id        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserId    primitive.ObjectID   `bson:"user_id" json:"user_id"`
	Products  []primitive.ObjectID `bson:"product_ids" json:"product_ids"`
	Discounts []Price              `bson:"discounts" json:"discounts"`
	Total     Price                `bson:"total" json:"total"`
}
