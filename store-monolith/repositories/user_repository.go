package repositories

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/manoamaro/microservice-store/monolith/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection string = "Users"

type UsersRepository struct {
	context    context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUsersRepository(mongoDB *mongo.Database) *UsersRepository {
	return &UsersRepository{
		context:    context.Background(),
		db:         mongoDB,
		collection: mongoDB.Collection(usersCollection),
	}
}

func (d *UsersRepository) Login(email string, password string) (*models.User, error) {
	res := d.collection.FindOne(d.context, bson.M{
		"email": bson.M{
			"$eq": email,
		},
	})

	result := &models.User{}
	if err := res.Decode(result); err != nil {
		return nil, err
	} else {
		if checkPassword(*result, password) {
			return result, nil
		} else {
			return nil, fmt.Errorf("invalid username and/or password")
		}
	}
}

func checkPassword(user models.User, plainPassword string) bool {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s+%s", plainPassword, user.Salt)))
	password := fmt.Sprintf("%x", h.Sum(nil))
	return user.Password == password
}
