package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const PRODUCTS_COLLECTION string = "Products"

type Product struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"title"`
	Description string             `bson:"description"`
	Images      []string           `bson:"images"`
	Reviews     []Review           `bson:"reviews"`
	Deleted     bool               `bson:"deleted"`
}

type Products interface {
	ListProducts() ([]Product, error)
	GetProduct(id primitive.ObjectID) (Product, error)
	DeleteProduct(id primitive.ObjectID) error
	InsertProduct(product Product) (string, error)
}
