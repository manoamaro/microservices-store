package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const USERS_COLLECTION = "Users"

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName string             `bson:"full_name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
	Salt     string             `bson:"salt" json:"-"`
	Deleted  bool               `bson:"deleted" json:"deleted"`
}

var UserEmailIndex = mongo.IndexModel{
	Keys:    bson.M{"email": 1},
	Options: options.Index().SetUnique(true),
}
