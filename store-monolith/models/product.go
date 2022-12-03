package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Prices      []Price            `bson:"prices" json:"prices"`
	Images      []string           `bson:"images" json:"images"`
	Reviews     []Review           `bson:"reviews" json:"reviews"`
	Inventory   int                `bson:"inventory" json:"inventory"`
	Deleted     bool               `bson:"deleted" json:"deleted"`
}

func (p Product) GetPrice(currency string) *Price {
	for _, v := range p.Prices {
		if v.Currency == currency {
			return &v
		}
	}
	return nil
}

type Price struct {
	Currency string `bson:"currency" json:"currency"`
	Price    int    `bson:"price" json:"price"`
}

type Review struct {
	UserId  primitive.ObjectID `bson:"user_id" json:"user_id"`
	Rating  int                `bson:"rating" json:"rating"`
	Comment string             `bson:"comment" json:"comment"`
}
