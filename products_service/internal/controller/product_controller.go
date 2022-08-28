package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manoamaro.github.com/commons/pkg/helpers"
	"manoamaro.github.com/products_service/internal/repository"
	"net/http"
)

type ProductController struct {
	ProductsRepository repository.ProductsRepository
}

func NewProductController(r *gin.Engine, productsRepository repository.ProductsRepository) *ProductController {
	controller := &ProductController{
		productsRepository,
	}
	publicGroup := r.Group("/public")
	{
		publicGroup.GET("/", controller.listProductsHandler)
		publicGroup.GET("/:id", controller.getProductHandler)
	}

	return controller
}

func (p *ProductController) listProductsHandler(c *gin.Context) {
	if products, err := p.ProductsRepository.ListProducts(); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, products)
	}
}

func (p *ProductController) getProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helpers.BadRequest(err, c)
		return
	}
	if product, err := p.ProductsRepository.GetProduct(objectID); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, product)
	}
}
