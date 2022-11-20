package repositories

import (
	"context"

	"github.com/manoamaro/microservice-store/monolith/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const productsCollection string = "Products"

type ProductsRepository struct {
	context    context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

func NewProductsRepository(db *mongo.Database) *ProductsRepository {
	return &ProductsRepository{
		context:    context.Background(),
		db:         db,
		collection: db.Collection(productsCollection),
	}
}

func (d *ProductsRepository) ListProducts() ([]models.Product, error) {
	cursor, err := d.collection.Find(d.context, bson.D{})
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

func (d *ProductsRepository) SearchProducts(query string) ([]models.Product, error) {
	cursor, err := d.collection.Find(d.context, bson.M{
		"$text": bson.M{
			"$search": query,
		},
	})
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

func (d *ProductsRepository) GetProduct(id primitive.ObjectID) (*models.Product, error) {
	res := d.collection.FindOne(d.context, bson.M{"_id": id})
	result := &models.Product{}
	if err := res.Decode(result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (d *ProductsRepository) DeleteProduct(id primitive.ObjectID) (bool, error) {
	if res, err := d.collection.DeleteOne(d.context, bson.M{"_id": id}); err != nil {
		return false, err
	} else {
		return res.DeletedCount > 0, nil
	}
}

func (d *ProductsRepository) InsertProduct(product models.Product) (*models.Product, error) {
	if res, err := d.collection.InsertOne(d.context, product); err != nil {
		return nil, err
	} else {
		newProduct := product
		newProduct.Id = res.InsertedID.(primitive.ObjectID)
		return &newProduct, nil
	}
}

func (d *ProductsRepository) UpdateProduct(product models.Product) (bool, error) {
	if res, err := d.collection.ReplaceOne(d.context, bson.M{"_id": product.Id}, product); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}
}

func (d *ProductsRepository) CreateReview(hexId string, userHexId string, stars int, comment string) (*models.Review, error) {
	if productId, err := primitive.ObjectIDFromHex(hexId); err != nil {
		return nil, err
	} else if userId, err := primitive.ObjectIDFromHex(userHexId); err != nil {
		return nil, err
	} else if product, err := d.GetProduct(productId); err != nil {
		return nil, err
	} else {
		review := &models.Review{
			UserId:  userId,
			Rating:  stars,
			Comment: comment,
		}

		product.Reviews = append(product.Reviews, *review)
		if _, err := d.UpdateProduct(*product); err != nil {
			return nil, err
		} else {
			return review, nil
		}
	}
}
