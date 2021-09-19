package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const USERS_COLLECTION = "Users"

type User struct {
    Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    FullName string             `bson:"full_name" json:"full_name"`
    Email    string             `bson:"email" json:"email"`
    Password string             `bson:"password" json:"password"`
    Salt     string             `bson:"salt" json:"salt"`
    Deleted  bool               `bson:"deleted" json:"deleted"`
}
