package repositories

import (
	"context"
	"fmt"

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

func (r *CartsRepository) GetOrCreateCart(userId primitive.ObjectID, currency string) (*models.Cart, error) {
	filter := bson.D{{"user_id", userId}}
	update := bson.D{
		{"$setOnInsert", bson.D{{"currency", currency}}},
	}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	if res := r.collection.FindOneAndUpdate(r.context, filter, update, options); res.Err() != nil {
		return nil, res.Err()
	} else {
		cart := &models.Cart{}
		if err := res.Decode(cart); err != nil {
			return nil, err
		} else {
			return cart, nil
		}
	}
}

func (r *CartsRepository) AddProduct(userHexId string, product models.Product, quantity int) (*models.Cart, error) {
	if userId, err := primitive.ObjectIDFromHex(userHexId); err != nil {
		return nil, err
	} else {
		cart, err := r.GetOrCreateCart(userId, "EUR")
		if err != nil {
			return nil, err
		}

		productPrice := Find(product.Prices, func(t models.Price) bool {
			return t.Currency == cart.Currency
		})

		if productPrice == nil {
			return nil, fmt.Errorf("products does not contain price in %s", cart.Currency)
		}

		filter := bson.D{{"user_id", userId}}
		update := bson.D{
			{"$addToSet", bson.D{
				{"products", bson.D{
					{"_id", product.Id},
					{"price", productPrice.Price},
					{"quantity", quantity},
				}},
			}},
			{"$set", bson.D{{
				"total", bson.D{{
					"$reduce", bson.D{
						{"input", "$products"},
						{"initialValue", 0},
						{"in", bson.D{
							{"$sum", bson.A{
								"$$value",
								bson.D{
									{"$multiply", bson.A{
										"$products.price", "$products.quantity",
									}},
								}},
							}},
						},
					}},
				}}}}}

		options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

		res := r.collection.FindOneAndUpdate(r.context, filter, update, options)

		if res.Err() != nil {
			return nil, res.Err()
		} else if err := res.Decode(cart); err != nil {
			return nil, err
		} else {
			return cart, nil
		}
	}
}

func Find[T interface{}](i []T, f func(T) bool) *T {
	for _, v := range i {
		if f(v) {
			return &v
		}
	}
	return nil
}
