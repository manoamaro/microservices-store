package mock

import (
	"github.com/manoamaro/microservices-store/products_service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/maps"
)

type MockProductsRepository struct {
	products map[string]models.Product
}

func NewMockProductsRepository() *MockProductsRepository {
	return &MockProductsRepository{products: make(map[string]models.Product)}
}

func (m *MockProductsRepository) ListProducts() ([]models.Product, error) {
	return maps.Values(m.products), nil
}

func (m *MockProductsRepository) GetProduct(id primitive.ObjectID) (*models.Product, error) {
	v := m.products[id.String()]
	return &v, nil
}

func (m *MockProductsRepository) DeleteProduct(id primitive.ObjectID) (bool, error) {
	delete(m.products, id.String())
	return true, nil
}

func (m *MockProductsRepository) InsertProduct(product models.Product) (*models.Product, error) {
	product.Id = primitive.NewObjectID()
	m.products[product.Id.String()] = product
	return &product, nil
}

func (m *MockProductsRepository) UpdateProduct(product models.Product) (bool, error) {
	m.products[product.Id.String()] = product
	return true, nil
}
