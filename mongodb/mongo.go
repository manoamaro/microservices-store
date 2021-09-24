package mongodb

import (
	"context"
	"crypto/sha512"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"time"
)

type MongoDB struct {
	Client   *mongo.Client
	Ctx      context.Context
	Database string
}

var Random = rand.New(rand.NewSource(time.Now().UnixNano()))

func ConnectMongoDB(url string, database string) *MongoDB {
	ctx := context.Background()
	if client, err := mongo.NewClient(options.Client().ApplyURI(url)); err != nil {
		return nil
	} else if err := client.Connect(ctx); err != nil {
		return nil
	} else {
		db := &MongoDB{Client: client, Ctx: ctx, Database: database}
		return db
	}
}

func (db *MongoDB) DisconnectMongoDB() error {
	return db.Client.Disconnect(db.Ctx)
}

func (db *MongoDB) Collection(collection string) *mongo.Collection {
	return db.Client.Database(db.Database).Collection(collection)
}

func Hash(value string) string {
	sha := sha512.New()
	sha.Write([]byte(value))
	return fmt.Sprintf("%x", sha.Sum(nil))
}
