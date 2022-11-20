package repositories

import (
	"context"

	"github.com/manoamaro/microservice-store/monolith/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CartsRepository struct {
	context    context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

func NewCartsRepository(mongoDB *mongo.Database) *CartsRepository {
	return &CartsRepository{
		context:    context.Background(),
		db:         mongoDB,
		collection: mongoDB.Collection("Carts"),
	}
}

func (r *CartsRepository) GetCartForUser(userHexId string) (*models.Cart, error) {
	if userId, err := primitive.ObjectIDFromHex(userHexId); err != nil {
		return nil, err
	} else {
		cart := &models.Cart{}

		res := r.collection.FindOne(r.context, bson.M{
			"user_id": bson.M{
				"$eq": userId,
			},
		})

		if err := res.Decode(cart); err != nil {
			return nil, err
		} else {
			return cart, nil
		}
	}
}

func (r *CartsRepository) AddProduct(userHexId string, product models.Product) (*models.Cart, error) {
	if userId, err := primitive.ObjectIDFromHex(userHexId); err != nil {
		return nil, err
	} else {
		productPrice := product.Prices[0]

		filter := bson.D{{"user_id", userId}}
		update := bson.D{
			{"$addToSet", bson.D{{"product_ids", product.Id}}},
			{"$set", bson.D{{"total.currency", productPrice.Currency}}},
			{"$inc", bson.D{{"total.price", productPrice.Price}}},
		}

		options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
		res := r.collection.FindOneAndUpdate(r.context, filter, update, options)

		cart := &models.Cart{}

		if res.Err() != nil {
			return nil, res.Err()
		} else if err := res.Decode(cart); err != nil {
			return nil, err
		} else {
			return cart, nil
		}
	}
}
