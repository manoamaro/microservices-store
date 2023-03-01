package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/manoamaro/microservices-store/products_service/internal"
	"github.com/manoamaro/microservices-store/products_service/internal/mock"
	"github.com/manoamaro/microservices-store/products_service/internal/models"
)

var productsRepository = mock.NewProductsRepository()

var application = internal.Application{
	ProductsRepository: productsRepository,
	AuthService:        &mock.AuthService{},
}

var router = application.SetupRoutes()

type ProductResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Prices      []struct {
		Id       string  `json:"id"`
		Currency string  `json:"currency"`
		Price    float64 `json:"price"`
	} `json:"prices"`
	Images  []string `json:"images"`
	Reviews []struct {
		Author  string `json:"author"`
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
}

func TestProductListEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/public/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}

func TestProductListWithItems(t *testing.T) {
	product, _ := productsRepository.InsertProduct(models.Product{Name: "Product1", Description: "DescriptionProduct1"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/public/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response []ProductResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, len(response), 1)
	assert.Equal(t, response[0].Id, product.Id.Hex())
	assert.Equal(t, response[0].Name, "Product1")
	assert.Equal(t, response[0].Description, "DescriptionProduct1")
	assert.Equal(t, len(response[0].Images), 0)
	assert.Equal(t, len(response[0].Prices), 0)
	assert.Equal(t, len(response[0].Reviews), 0)
}

func TestProductGet(t *testing.T) {
	product, _ := productsRepository.InsertProduct(models.Product{Name: "Product1", Description: "DescriptionProduct1"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/public/"+product.Id.Hex(), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response ProductResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, response.Id, product.Id.Hex())
	assert.Equal(t, response.Name, "Product1")
	assert.Equal(t, response.Description, "DescriptionProduct1")
	assert.Equal(t, len(response.Images), 0)
	assert.Equal(t, len(response.Prices), 0)
	assert.Equal(t, len(response.Reviews), 0)
}

func TestProductAdminCreate(t *testing.T) {

	request := models.CreateProductRequest{
		Name:        "Product admin Create 1",
		Description: "Description Product admin Create 1",
		Price: models.PriceRequest{
			Currency: "EUR",
			Price:    12.58,
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/admin/create", bytes.NewReader(requestBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusCreated)
}
