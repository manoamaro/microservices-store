package repositories

import (
	"fmt"

	"github.com/manoamaro/microservice-store/monolith/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartsRepository struct {
	r                 *Repository[models.Cart]
	productRepository *ProductsRepository
}

func NewCartsRepository(mongoDB *mongo.Database, productRepository *ProductsRepository) *CartsRepository {
	return &CartsRepository{
		r:                 (*Repository[models.Cart])(NewRepository[models.Cart](mongoDB, "Carts")),
		productRepository: productRepository,
	}
}

func (r *CartsRepository) GetCartForUser(userHexId string) (*models.Cart, error) {
	if userId, err := primitive.ObjectIDFromHex(userHexId); err != nil {
		return nil, err
	} else {
		return r.r.Find(bson.M{"user_id": userId})
	}
}

func (r *CartsRepository) GetOrCreateCart(userId primitive.ObjectID, currency string) (*models.Cart, error) {
	filter := bson.M{"user_id": userId}
	update := bson.M{
		"$setOnInsert": bson.M{"currency": currency},
	}

	return r.r.FindOneAndUpdate(filter, update, true)
}

func (r *CartsRepository) AddProduct(userHexId string, product models.Product, quantity int) (*models.Cart, error) {
	if userId, err := primitive.ObjectIDFromHex(userHexId); err != nil {
		return nil, err
	} else {
		cart, err := r.GetOrCreateCart(userId, "EUR")
		if err != nil {
			return nil, err
		}

		price := product.GetPrice(cart.Currency)
		if price == nil {
			return nil, fmt.Errorf("product does not have price in %s", cart.Currency)
		}

		if inventoryUpdated, err := r.productRepository.DicreaseInventory(product.Id.Hex(), quantity); err != nil {
			return nil, err
		} else if !inventoryUpdated {
			return nil, fmt.Errorf("inventory could not get updated")
		}

		filter := bson.M{"_id": cart.Id}

		setAddProduct := bson.M{
			"$set": bson.M{
				"products": bson.M{
					"$cond": bson.A{
						bson.M{"$in": bson.A{product.Id, "$products.id"}},
						bson.M{
							"$map": bson.D{
								{"input", "$products"},
								{"in", bson.M{
									"$cond": bson.A{
										bson.M{"$gt": bson.A{bson.M{"$add": bson.A{"$$this.quantity", quantity}}, 0}},
										bson.M{
											"$mergeObjects": bson.A{
												"$$this",
												bson.M{
													"$cond": bson.A{
														bson.M{"$eq": bson.A{"$$this.id", product.Id}},
														bson.D{
															{"quantity", bson.M{"$add": bson.A{"$$this.quantity", quantity}}},
															{"value", price.Price},
														},
														bson.M{},
													},
												},
											},
										},
										nil,
									},
								}},
							},
						},
						bson.M{
							"$cond": bson.A{
								bson.M{"$gt": bson.A{quantity, 0}},
								bson.M{
									"$concatArrays": bson.A{"$products", bson.A{
										bson.D{
											{"id", product.Id},
											{"value", price.Price},
											{"quantity", quantity},
										},
									}},
								},
								"$products",
							},
						},
					},
				},
			},
		}

		setTotal := bson.M{
			"$set": bson.M{
				"total": bson.M{
					"$reduce": bson.D{
						{"input", "$products"},
						{"initialValue", 0},
						{"in", bson.M{
							"$add": bson.A{
								"$$value",
								bson.M{
									"$multiply": bson.A{"$$this.quantity", "$$this.value"},
								},
							},
						}},
					},
				},
			},
		}

		update := bson.A{
			setAddProduct,
			bson.M{
				"$set": bson.M{
					"products": bson.M{
						"$filter": bson.D{
							{"input", "$products"},
							{"cond", bson.M{"$ne": bson.A{"$$this", nil}}},
						},
					},
				},
			},
			setTotal,
			bson.M{"$set": bson.M{"updatedAt": "$NOW"}},
		}

		return r.r.FindOneAndUpdate(filter, update, true)
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
