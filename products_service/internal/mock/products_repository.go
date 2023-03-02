package mock

import (
	"github.com/manoamaro/microservices-store/products_service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/maps"
)

type ProductsRepository struct {
	products map[string]models.Product
}

func NewProductsRepository() *ProductsRepository {
	return &ProductsRepository{products: make(map[string]models.Product)}
}

func (d *ProductsRepository) ListProducts() ([]models.Product, error) {
	return maps.Values(d.products), nil
}

func (d *ProductsRepository) SearchProducts(query string) ([]models.Product, error) {
	return make([]models.Product, 0), nil
}

func (d *ProductsRepository) GetProduct(id primitive.ObjectID) (*models.Product, error) {
	v := d.products[id.String()]
	return &v, nil
}

func (d *ProductsRepository) DeleteProduct(id primitive.ObjectID) (bool, error) {
	delete(d.products, id.String())
	return true, nil
}

func (d *ProductsRepository) InsertProduct(product models.Product) (*models.Product, error) {
	product.Id = primitive.NewObjectID()
	d.products[product.Id.String()] = product
	return &product, nil
}

func (d *ProductsRepository) UpdateProduct(product models.Product) (bool, error) {
	d.products[product.Id.String()] = product
	return true, nil
}

func (d *ProductsRepository) CreateReview(productHexId string, userId string, rating int, comment string) (*models.Review, error) {
	return nil, nil
}
