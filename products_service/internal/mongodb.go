package internal

import (
	"manoamaro.github.com/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manoamaro.github.com/products_service/internal/models"
)

type DB struct {
	*mongodb.MongoDB
}

const DATABASE string = "products"

func (db *DB) ListProducts() ([]models.Product, error) {
	if cur, err := db.Collection(models.PRODUCTS_COLLECTION).Find(db.Ctx, bson.D{}); err != nil {
		return nil, err
	} else {
		var result []models.Product
		if err = cur.All(db.Ctx, &result); err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (db *DB) InsertProduct(product models.Product) (*models.Product, error) {
	if res, err := db.Collection(models.PRODUCTS_COLLECTION).InsertOne(db.Ctx, product); err != nil {
		return nil, err
	} else {
		product.Id = res.InsertedID.(primitive.ObjectID)
		return &product, nil
	}
}

func (db *DB) UpdateProduct(id primitive.ObjectID, product models.Product) (bool, error) {
	if res, err := db.Collection(models.PRODUCTS_COLLECTION).ReplaceOne(db.Ctx, bson.M{"_id": id}, product); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}
}

func (db *DB) FetchProduct(id primitive.ObjectID) (*models.Product, error) {
	res := db.Collection(models.PRODUCTS_COLLECTION).FindOne(db.Ctx, bson.M{"_id": id})
	result := &models.Product{}
	if err := res.Decode(result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (db *DB) DeleteProduct(id primitive.ObjectID) (bool, error) {
	if res, err := db.Collection(models.PRODUCTS_COLLECTION).DeleteOne(db.Ctx, bson.M{"_id": id}); err != nil {
		return false, err
	} else {
		return res.DeletedCount > 0, nil
	}
}
