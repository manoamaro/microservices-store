package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository[T any] struct {
	context    context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

func NewRepository[T any](mongoDB *mongo.Database, collection string) *Repository[T] {
	return &Repository[T]{
		context:    context.Background(),
		db:         mongoDB,
		collection: mongoDB.Collection(collection),
	}
}

func (r *Repository[T]) Find(q interface{}) (*T, error) {
	res := r.collection.FindOne(r.context, q)
	if res.Err() != nil {
		return nil, res.Err()
	} else {
		result := new(T)
		if err := res.Decode(result); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

func (r *Repository[T]) FindById(id primitive.ObjectID) (*T, error) {
	return r.Find(bson.M{"_id": id})
}

func (r *Repository[T]) Query(q bson.D) ([]T, error) {
	if cur, err := r.collection.Find(r.context, q); err != nil {
		return nil, err
	} else {
		result := new([]T)
		if err := cur.All(r.context, result); err != nil {
			return nil, err
		} else {
			return *result, nil
		}
	}
}

func (r *Repository[T]) List() ([]T, error) {
	return r.Query(bson.D{})
}

func (r *Repository[T]) Replace(q bson.D, doc T) (bool, error) {
	if res, err := r.collection.ReplaceOne(r.context, q, doc); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0, nil
	}
}

func (r *Repository[T]) FindOneAndUpdate(q interface{}, u interface{}, upsert bool) (*T, error) {
	options := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(upsert)
	res := r.collection.FindOneAndUpdate(r.context, q, u, options)
	if res.Err() != nil {
		return nil, res.Err()
	} else {
		updated := new(T)
		if err := res.Decode(updated); err != nil {
			return nil, err
		} else {
			return updated, nil
		}
	}
}

func (r *Repository[T]) UpdateOne(q interface{}, u interface{}, upsert bool) (bool, error) {
	options := options.Update().SetUpsert(upsert)
	if res, err := r.collection.UpdateOne(r.context, q, u, options); err != nil {
		return false, err
	} else {
		return res.ModifiedCount > 0 || res.UpsertedCount > 0, nil
	}
}

func (r *Repository[T]) Delete(q interface{}) (bool, error) {
	if res, err := r.collection.DeleteOne(r.context, q); err != nil {
		return false, err
	} else {
		return res.DeletedCount > 0, nil
	}
}

func (r *Repository[T]) Insert(doc T) (*T, error) {
	if res, err := r.collection.InsertOne(r.context, doc); err != nil {
		return nil, err
	} else {
		return r.FindById(res.InsertedID.(primitive.ObjectID))
	}
}
