package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"manoamaro.github.com/products_service/internal"
)

const ProductsCollection string = "Products"

type Product struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Images      []string           `bson:"images" json:"images"`
	Reviews     []Review           `bson:"reviews" json:"reviews"`
	Deleted     bool               `bson:"deleted" json:"deleted"`
}

type ProductsRepository interface {
	ListProducts() ([]Product, error)
	GetProduct(id primitive.ObjectID) (*Product, error)
	DeleteProduct(id primitive.ObjectID) (bool, error)
	InsertProduct(product Product) (*Product, error)
	UpdateProduct(product Product) (bool, error)
}

type DefaultProductsRepository struct {
	context context.Context
	client  *mongo.Client
}

func NewProductsRepository() ProductsRepository {
	uri := internal.GetEnv("MONGO_URL", "mongodb://test:test@localhost:27017/?maxPoolSize=20&w=majority")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return &DefaultProductsRepository{
		context: context.Background(),
		client:  client,
	}
}

func (d *DefaultProductsRepository) ListProducts() ([]Product, error) {
	cursor, err := d.client.Database(internal.ProductsServiceDatabase).Collection(ProductsCollection).Find(d.context, bson.D{})
	if err != nil {
		return nil, err
	}
	var result []Product
	err = cursor.All(d.context, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DefaultProductsRepository) GetProduct(id primitive.ObjectID) (*Product, error) {
	res := d.client.Database(internal.ProductsServiceDatabase).Collection(ProductsCollection).FindOne(d.context, bson.M{"_id": id})
	result := &Product{}
	if err := res.Decode(result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (d *DefaultProductsRepository) DeleteProduct(id primitive.ObjectID) (bool, error) {
	if res, err := d.client.Database(internal.ProductsServiceDatabase).Collection(ProductsCollection).DeleteOne(d.context, bson.M{"_id": id}); err != nil {
		return false, err
	} else {
		return res.DeletedCount > 0, nil
	}
}

func (d *DefaultProductsRepository) InsertProduct(product Product) (*Product, error) {
	if res, err := d.client.Database(internal.ProductsServiceDatabase).Collection(ProductsCollection).InsertOne(d.context, product); err != nil {
		return nil, err
	} else {
		newProduct := product
		newProduct.Id = res.InsertedID.(primitive.ObjectID)
		return &newProduct, nil
	}
}

func (d *DefaultProductsRepository) UpdateProduct(product Product) (bool, error) {
	if res, err := d.client.Database(internal.ProductsServiceDatabase).Collection(ProductsCollection).ReplaceOne(d.context, bson.M{"_id": product.Id}, product); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}

}
