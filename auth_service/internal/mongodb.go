package internal

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"manoamaro.github.com/auth_service/models"
	"math/rand"
	"strconv"
	"time"
)

type MongoDB struct {
	client *mongo.Client
	ctx    context.Context
}

const DATABASE string = "users"

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func ConnectMongoDB(url string) *MongoDB {
	ctx := context.Background()
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	FailOnError(err)
	err = client.Connect(ctx)
	FailOnError(err)

	db := &MongoDB{client: client, ctx: ctx}

	if _, err := db.collection().Indexes().CreateOne(db.ctx, models.UserEmailIndex); err != nil {
		log.Println(err)
	}

	return db
}

func (db *MongoDB) DisconnectMongoDB() error {
	return db.client.Disconnect(db.ctx)
}

func (db *MongoDB) collection() *mongo.Collection {
	return db.client.Database(DATABASE).Collection(models.USERS_COLLECTION)
}

func hash(value string) string {
	sha := sha512.New()
	sha.Write([]byte(value))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func (db *MongoDB) LoginUser(email string, plainPassword string) (*models.User, error) {
	if res := db.collection().FindOne(db.ctx, bson.M{"email": email}); res.Err() != nil {
		return nil, res.Err()
	} else {
		result := &models.User{}
		if err := res.Decode(result); err != nil {
			return nil, err
		}
		pwd := hash(plainPassword + result.Salt)
		if pwd != result.Password {
			return nil, errors.New("invalid password")
		}

		return result, nil
	}
}

func (db *MongoDB) CreateUser(user models.User, plainPassword string) (*models.User, error) {
	salt := hash(strconv.Itoa(random.Int()))
	pwd := hash(plainPassword + salt)
	user.Salt = salt
	user.Password = pwd

	if res, err := db.collection().InsertOne(db.ctx, user); err != nil {
		return nil, err
	} else {
		user.Id = res.InsertedID.(primitive.ObjectID)
		return &user, nil
	}
}

func (db *MongoDB) UpdateProduct(id primitive.ObjectID, user models.User) (bool, error) {
	if res, err := db.collection().ReplaceOne(db.ctx, bson.M{"_id": id}, user); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}
}

func (db *MongoDB) FetchProduct(id primitive.ObjectID) (*models.User, error) {
	res := db.collection().FindOne(db.ctx, bson.M{"_id": id})
	result := &models.User{}
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
