package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const PRODUCTS_COLLECTION string = "Products"

type Product struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Images      []string           `bson:"images" json:"images"`
	Reviews     []Review           `bson:"reviews" json:"reviews"`
	Deleted     bool               `bson:"deleted" json:"deleted"`
}

type Products interface {
	ListProducts() ([]Product, error)
	GetProduct(id primitive.ObjectID) (Product, error)
	DeleteProduct(id primitive.ObjectID) error
	InsertProduct(product Product) (string, error)
}
