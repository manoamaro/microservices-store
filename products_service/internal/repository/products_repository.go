package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"manoamaro.github.com/commons"
	"manoamaro.github.com/products_service/internal/models"
)

const ProductsCollection string = "Products"
const ProductsServiceDatabase = "ProductsService"

type ProductsRepository interface {
	ListProducts() ([]models.Product, error)
	GetProduct(id primitive.ObjectID) (*models.Product, error)
	DeleteProduct(id primitive.ObjectID) (bool, error)
	InsertProduct(product models.Product) (*models.Product, error)
	UpdateProduct(product models.Product) (bool, error)
}

type DefaultProductsRepository struct {
	context context.Context
	client  *mongo.Client
	db      *mongo.Database
}

func NewDefaultProductsRepository() ProductsRepository {
	uri := commons.GetEnv("MONGO_URL", "mongodb://test:test@localhost:27017/?maxPoolSize=20&w=majority")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return &DefaultProductsRepository{
		context: context.Background(),
		client:  client,
		db:      client.Database(ProductsServiceDatabase),
	}
}

func (d *DefaultProductsRepository) ListProducts() ([]models.Product, error) {
	cursor, err := d.client.Database(ProductsServiceDatabase).Collection(ProductsCollection).Find(d.context, bson.D{})
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
