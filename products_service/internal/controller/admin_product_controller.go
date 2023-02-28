package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/products_service/internal/models"
	"github.com/manoamaro/microservices-store/products_service/internal/repository"
	"net/http"
)

type AdminProductController struct {
	authService infra.AuthService
	ProductController
}

func NewAdminProductController(r *gin.Engine, productsRepository repository.ProductsRepository, authService infra.AuthService) *AdminProductController {
	productsController := ProductController{productsRepository}
	controller := &AdminProductController{
		authService,
		productsController,
	}

	mgmtGroup := r.Group("/admin")
	{
		mgmtGroup.Use(helpers.AuthMiddleware(authService, "products_admin"))
		mgmtGroup.GET("/list", controller.listProductsHandler)
		mgmtGroup.POST("/create", controller.postProductsHandler)
	}
	return controller
}

func (a *AdminProductController) postProductsHandler(c *gin.Context) {
	newProduct := models.Product{}
	if err := c.BindJSON(&newProduct); err != nil {
		helpers.BadRequest(err, c)
	} else if savedProduct, err := a.ProductsRepository.InsertProduct(newProduct); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(http.StatusCreated, savedProduct)
	}
}
