package internal

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manoamaro.github.com/auth_service/models"
	"manoamaro.github.com/mongodb"
	"strconv"
)

const DATABASE string = "users"

type DB struct {
	*mongodb.MongoDB
}

func (db *DB) LoginUser(email string, plainPassword string) (*models.User, error) {
	if res := db.Collection(models.USERS_COLLECTION).FindOne(db.Ctx, bson.M{"email": email}); res.Err() != nil {
		return nil, res.Err()
	} else {
		result := &models.User{}
		if err := res.Decode(result); err != nil {
			return nil, err
		}
		pwd := mongodb.Hash(plainPassword + result.Salt)
		if pwd != result.Password {
			return nil, errors.New("invalid password")
		}

		return result, nil
	}
}

func (db *DB) CreateUser(user models.User, plainPassword string) (*models.User, error) {
	salt := mongodb.Hash(strconv.Itoa(mongodb.Random.Int()))
	pwd := mongodb.Hash(plainPassword + salt)
	user.Salt = salt
	user.Password = pwd

	if res, err := db.Collection(models.USERS_COLLECTION).InsertOne(db.Ctx, user); err != nil {
		return nil, err
	} else {
		user.Id = res.InsertedID.(primitive.ObjectID)
		return &user, nil
	}
}

func (db *DB) UpdateProduct(id primitive.ObjectID, user models.User) (bool, error) {
	if res, err := db.Collection(models.USERS_COLLECTION).ReplaceOne(db.Ctx, bson.M{"_id": id}, user); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}
}

func (db *DB) FetchProduct(id primitive.ObjectID) (*models.User, error) {
	res := db.Collection(models.USERS_COLLECTION).FindOne(db.Ctx, bson.M{"_id": id})
	result := &models.User{}
	if err := res.Decode(result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (db *DB) DeleteProduct(id primitive.ObjectID) (bool, error) {
	if res, err := db.Collection(models.USERS_COLLECTION).DeleteOne(db.Ctx, bson.M{"_id": id}); err != nil {
		return false, err
	} else {
		return res.DeletedCount > 0, nil
	}
}
