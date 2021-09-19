package internal

import (
	"context"

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
	ctx := context.Background()
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	FailOnError(err)
	err = client.Connect(ctx)
	FailOnError(err)
	DB = &MongoDB{client: client, ctx: ctx}
	return client
}

func DisconnectMongoDB() error {
	return DB.client.Disconnect(DB.ctx)
}

func (db *MongoDB) collection() *mongo.Collection {
	return db.client.Database(DATABASE).Collection(models.PRODUCTS_COLLECTION)
}

func (db *MongoDB) ListProducts() ([]models.Product, error) {
	if cur, err := db.collection().Find(db.ctx, bson.D{}); err != nil {
		return nil, err
	} else {
		var result []models.Product
		if err = cur.All(db.ctx, &result); err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (db *MongoDB) InsertProduct(product models.Product) (*models.Product, error) {
	if res, err := db.collection().InsertOne(db.ctx, product); err != nil {
		return nil, err
	} else {
		product.Id = res.InsertedID.(primitive.ObjectID)
		return &product, nil
	}
}

func (db *MongoDB) UpdateProduct(id primitive.ObjectID, product models.Product) (bool, error) {
	if res, err := db.collection().ReplaceOne(db.ctx, bson.M{"_id": id}, product); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}
}

func (db *MongoDB) FetchProduct(id primitive.ObjectID) (*models.Product, error) {
	res := db.collection().FindOne(db.ctx, bson.M{"_id": id})
	result := &models.Product{}
	if err := res.Decode(result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (db *MongoDB) DeleteProduct(id primitive.ObjectID) (bool, error) {
	if res, err := db.collection().DeleteOne(db.ctx, bson.M{"_id": id}); err != nil {
		return false, err
	} else {
		return res.DeletedCount > 0, nil
	}
}
