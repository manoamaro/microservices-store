package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Prices      []Price            `bson:"prices"`
	Images      []string           `bson:"images"`
	Reviews     []Review           `bson:"reviews"`
	Deleted     bool               `bson:"deleted"`
}

type Price struct {
	Currency string `bson:"currency"`
	Value    int    `bson:"price"`
}

type Review struct {
	UserId  string `bson:"user_id"`
	Author  string `bson:"author" `
	Rating  int    `bson:"rating"`
	Comment string `bson:"comment"`
}
