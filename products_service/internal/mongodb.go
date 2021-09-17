package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"manoamaro.github.com/products_service/internal/models"
)

type MongoDB struct {
	client *mongo.Client
	ctx    context.Context
}

var DB *MongoDB

const DATABASE string = "products"

func ConnectMongoDB(url string) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	FailOnError(err)
	err = client.Connect(ctx)
	FailOnError(err)
	DB = &MongoDB{client: client, ctx: ctx}
	return client
}

func (db *MongoDB) collection(name string) *mongo.Collection {
	return db.client.Database(DATABASE).Collection(name)
}

func (db *MongoDB) ListProducts() ([]models.Product, error) {
	cur, err := db.collection(models.PRODUCTS_COLLECTION).Find(db.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var result []models.Product
	err = cur.All(db.ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *MongoDB) InsertProduct(product models.Product) (*models.Product, error) {
	res, err := db.collection(models.PRODUCTS_COLLECTION).InsertOne(db.ctx, product)
	if err != nil {
		return nil, err
	}

	product.Id = res.InsertedID.(primitive.ObjectID)
	return &product, nil
}

func (db *MongoDB) FetchProduct(id primitive.ObjectID) (*models.Product, error) {
	res := db.collection(models.PRODUCTS_COLLECTION).FindOne(db.ctx, bson.M{"_id": id})
	result := &models.Product{}
	err := res.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
