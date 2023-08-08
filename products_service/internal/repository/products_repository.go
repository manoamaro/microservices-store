package repository

import (
	"context"

	"github.com/manoamaro/microservices-store/products_service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockery --name ProductsRepository --case=snake --output ../../test/mocks
type ProductsRepository interface {
	ListProducts() ([]models.Product, error)
	SearchProducts(query string) ([]models.Product, error)
	GetProduct(id primitive.ObjectID) (*models.Product, error)
	DeleteProduct(id primitive.ObjectID) (bool, error)
	InsertProduct(product models.Product) (*models.Product, error)
	UpdateProduct(product models.Product) (bool, error)
	AddImage(id primitive.ObjectID, imagePath string) (bool, error)
	DeleteImage(id primitive.ObjectID, imageId string) (bool, error)
	CreateReview(productHexId string, userId string, rating int, comment string) (*models.Review, error)
}

type DefaultProductsRepository struct {
	context context.Context
	col     *mongo.Collection
}

const ProductsCollection string = "Products"

func NewDefaultProductsRepository(db *mongo.Database) ProductsRepository {
	return &DefaultProductsRepository{
		context: context.Background(),
		col:     db.Collection(ProductsCollection),
	}
}

func (d *DefaultProductsRepository) ListProducts() ([]models.Product, error) {
	cursor, err := d.col.Find(d.context, bson.D{})
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

func (d *DefaultProductsRepository) SearchProducts(query string) ([]models.Product, error) {
	q := bson.D{{
		Key: "$text", Value: bson.M{
			"$search": query,
		},
	}}
	var result []models.Product
	if cur, err := d.col.Find(d.context, q); err != nil {
		return nil, err
	} else {
		if err := cur.All(d.context, result); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

func (d *DefaultProductsRepository) GetProduct(id primitive.ObjectID) (*models.Product, error) {
	res := d.col.FindOne(d.context, bson.M{"_id": id})
	result := &models.Product{}
	if err := res.Decode(result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (d *DefaultProductsRepository) DeleteProduct(id primitive.ObjectID) (bool, error) {
	if res, err := d.col.DeleteOne(d.context, bson.M{"_id": id}); err != nil {
		return false, err
	} else {
		return res.DeletedCount > 0, nil
	}
}

func (d *DefaultProductsRepository) InsertProduct(product models.Product) (*models.Product, error) {
	if res, err := d.col.InsertOne(d.context, product); err != nil {
		return nil, err
	} else {
		newProduct := product
		newProduct.Id = res.InsertedID.(primitive.ObjectID)
		return &newProduct, nil
	}
}

func (d *DefaultProductsRepository) UpdateProduct(product models.Product) (bool, error) {
	if res, err := d.col.ReplaceOne(d.context, bson.M{"_id": product.Id}, product); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}
}

func (d *DefaultProductsRepository) AddImage(id primitive.ObjectID, imagePath string) (bool, error) {
	update := primitive.A{
		primitive.M{
			"$set": primitive.M{
				"images": primitive.M{
					"$ifNull": primitive.A{
						primitive.M{"$concatArrays": primitive.A{"$images", primitive.A{imagePath}}},
						primitive.A{imagePath},
					}}},
		},
	}
	if res, err := d.col.UpdateByID(d.context, id, update); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}
}

func (d *DefaultProductsRepository) DeleteImage(id primitive.ObjectID, imageId string) (bool, error) {
	update := primitive.A{
		primitive.M{
			"$set": primitive.M{
				"images": primitive.M{
					"$filter": primitive.M{
						"input": "$images",
						"as":    "image",
						"cond":  primitive.M{"$ne": primitive.A{"$$image", imageId}},
					}}},
		},
	}

	if res, err := d.col.UpdateByID(d.context, id, update); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}
}

func (d *DefaultProductsRepository) CreateReview(productHexId string, userId string, rating int, comment string) (*models.Review, error) {
	if productId, err := primitive.ObjectIDFromHex(productHexId); err != nil {
		return nil, err
	} else if product, err := d.GetProduct(productId); err != nil {
		return nil, err
	} else {

		review := &models.Review{
			UserId:  userId,
			Author:  userId,
			Rating:  rating,
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
