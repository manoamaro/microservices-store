package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Salt     string             `bson:"salt" json:"salt"`
	Name     string             `bson:"name" json:"name"`
	Deleted  bool               `bson:"deleted" json:"deleted"`
}
