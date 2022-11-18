package repository

import (
	"context"
	"github.com/manoamaro/microservices-store/products_service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const ProductsCollection string = "Products"

type ProductsRepository interface {
	ListProducts() ([]models.Product, error)
	GetProduct(id primitive.ObjectID) (*models.Product, error)
	DeleteProduct(id primitive.ObjectID) (bool, error)
	InsertProduct(product models.Product) (*models.Product, error)
	UpdateProduct(product models.Product) (bool, error)
}

type DefaultProductsRepository struct {
	context context.Context
	db      *mongo.Database
}

func NewDefaultProductsRepository(db *mongo.Database) ProductsRepository {
	return &DefaultProductsRepository{
		context: context.Background(),
		db:      db,
	}
}

func (d *DefaultProductsRepository) ListProducts() ([]models.Product, error) {
	cursor, err := d.db.Collection(ProductsCollection).Find(d.context, bson.D{})
	if err != nil {
		return nil, err
	}
	var result []models.Product
	err = cursor.All(d.context, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DefaultProductsRepository) GetProduct(id primitive.ObjectID) (*models.Product, error) {
	res := d.db.Collection(ProductsCollection).FindOne(d.context, bson.M{"_id": id})
	result := &models.Product{}
	if err := res.Decode(result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (d *DefaultProductsRepository) DeleteProduct(id primitive.ObjectID) (bool, error) {
	if res, err := d.db.Collection(ProductsCollection).DeleteOne(d.context, bson.M{"_id": id}); err != nil {
		return false, err
	} else {
		return res.DeletedCount > 0, nil
	}
}

func (d *DefaultProductsRepository) InsertProduct(product models.Product) (*models.Product, error) {
	if res, err := d.db.Collection(ProductsCollection).InsertOne(d.context, product); err != nil {
		return nil, err
	} else {
		newProduct := product
		newProduct.Id = res.InsertedID.(primitive.ObjectID)
		return &newProduct, nil
	}
}

func (d *DefaultProductsRepository) UpdateProduct(product models.Product) (bool, error) {
	if res, err := d.db.Collection(ProductsCollection).ReplaceOne(d.context, bson.M{"_id": product.Id}, product); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}

}
