package models

type Review struct {
	Author  string `bson:"author"`
	Rating  int    `bson:"rating"`
	Comment string `bson:"comment"`
}
