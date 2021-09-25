package internal

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manoamaro.github.com/cart_service/models"
	"manoamaro.github.com/mongodb"
)

type DB struct {
	*mongodb.MongoDB
}

const DATABASE string = "cart"

func (db *DB) GetCartForUser(userId string) (*models.Cart, error) {
	if objId, err := primitive.ObjectIDFromHex(userId); err != nil {
		return nil, err
	} else if res := db.Collection(models.CARTS_COLLECTION).FindOne(db.Ctx, bson.M{"user_id": objId}); res.Err() != nil {
		return nil, err
	} else {
		result := &models.Cart{}
		if err = res.Decode(result); err != nil {
			return nil, err
		}
		return result, nil
	}
}
