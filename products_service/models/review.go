package models

type Review struct {
	Author  string `bson:"author" json:"author"`
	Rating  int    `bson:"rating" json:"rating"`
	Comment string `bson:"comment" json:"comment"`
}
