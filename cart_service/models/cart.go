package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const CARTS_COLLECTION = "carts"

type Cart struct {
	UserId    primitive.ObjectID `bson:"user_id" json:"userId"`
	ProductId primitive.ObjectID `bson:"product_id" json:"productId"`
	Quantity  int                `bson:"quantity" json:"quantity"`
}
