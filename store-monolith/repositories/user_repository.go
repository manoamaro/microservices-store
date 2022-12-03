package repositories

import (
	"crypto/sha256"
	"fmt"

	"github.com/manoamaro/microservice-store/monolith/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection string = "Users"

type UsersRepository struct {
	r *Repository[models.User]
}

func NewUsersRepository(mongoDB *mongo.Database) *UsersRepository {
	return &UsersRepository{
		r: (*Repository[models.User])(NewRepository[models.User](mongoDB, usersCollection)),
	}
}

func (d *UsersRepository) Login(email string, password string) (*models.User, error) {
	if user, err := d.r.Find(bson.M{"email": email}); err != nil {
		return nil, err
	} else if checkPassword(*user, password) {
		return user, nil
	} else {
		return nil, fmt.Errorf("invalid username and/or password")
	}
}

func checkPassword(user models.User, plainPassword string) bool {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s+%s", plainPassword, user.Salt)))
	password := fmt.Sprintf("%x", h.Sum(nil))
	return user.Password == password
}
