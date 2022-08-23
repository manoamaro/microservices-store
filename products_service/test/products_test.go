package test

import (
	"github.com/go-playground/assert/v2"
	"manoamaro.github.com/products_service/internal"
	"manoamaro.github.com/products_service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

var productsRepository = NewMockProductsRepository()

var application = internal.Application{
	ProductsRepository: productsRepository,
	AuthService:        &MockAuthService{},
}

var router = application.SetupRoutes()

func TestProductListEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/public/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}

func TestProductListWithItems(t *testing.T) {
	_, _ = productsRepository.InsertProduct(models.Product{Name: "Product1", Description: "DescriptionProduct1"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/public/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}
