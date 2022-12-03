package repositories

import (
	"github.com/manoamaro/microservice-store/monolith/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const productsCollection string = "Products"

type ProductsRepository struct {
	r Repository[models.Product]
}

func NewProductsRepository(db *mongo.Database) *ProductsRepository {
	return &ProductsRepository{
		r: Repository[models.Product](*NewRepository[models.Product](db, productsCollection)),
	}
}

func (d *ProductsRepository) ListProducts() ([]models.Product, error) {
	return d.r.List()
}

func (d *ProductsRepository) SearchProducts(query string) ([]models.Product, error) {
	q := bson.D{{
		Key: "$text", Value: bson.M{
			"$search": query,
		},
	}}
	return d.r.Query(q)
}

func (d *ProductsRepository) GetProduct(id primitive.ObjectID) (*models.Product, error) {
	return d.r.FindById(id)
}

func (d *ProductsRepository) DeleteProduct(id primitive.ObjectID) (bool, error) {
	return d.r.Delete(bson.M{"_id": id})
}

func (d *ProductsRepository) InsertProduct(product models.Product) (*models.Product, error) {
	return d.r.Insert(product)
}

func (d *ProductsRepository) UpdateProduct(product models.Product) (bool, error) {
	return d.r.Replace(bson.D{{Key: "_id", Value: product.Id}}, product)
}

func (d *ProductsRepository) CreateReview(hexId string, userHexId string, stars int, comment string) (*models.Review, error) {
	if productId, err := primitive.ObjectIDFromHex(hexId); err != nil {
		return nil, err
	} else if userId, err := primitive.ObjectIDFromHex(userHexId); err != nil {
		return nil, err
	} else if product, err := d.GetProduct(productId); err != nil {
		return nil, err
	} else {
		review := &models.Review{
			UserId:  userId,
			Rating:  stars,
			Comment: comment,
		}

		product.Reviews = append(product.Reviews, *review)
		if _, err := d.UpdateProduct(*product); err != nil {
			return nil, err
		} else {
			return review, nil
		}
	}
}

func (d *ProductsRepository) DicreaseInventory(hexId string, quantity int) (bool, error) {
	if productId, err := primitive.ObjectIDFromHex(hexId); err != nil {
		return false, err
	} else {
		filter := bson.D{
			{"_id", productId},
			{"inventory", bson.D{{"$gte", quantity}}},
		}

		update := bson.A{
			bson.D{{"$set", bson.D{{"inventory", bson.D{{"$subtract", bson.A{"$inventory", quantity}}}}}}},
			bson.D{{"$set", bson.D{{"updatedAt", "$$NOW"}}}},
		}

		return d.r.UpdateOne(filter, update, false)
	}
}
