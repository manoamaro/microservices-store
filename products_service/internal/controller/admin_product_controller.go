package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"manoamaro.github.com/commons/pkg/helpers"
	"manoamaro.github.com/commons/pkg/services"
	"manoamaro.github.com/products_service/internal/models"
	"manoamaro.github.com/products_service/internal/repository"
	"net/http"
)

type AdminProductController struct {
	authService services.AuthService
	ProductController
}

func NewAdminProductController(r *gin.Engine, productsRepository repository.ProductsRepository, authService services.AuthService) *AdminProductController {
	productsController := ProductController{productsRepository}
	controller := &AdminProductController{
		authService,
		productsController,
	}

	mgmtGroup := r.Group("/admin")
	{
		mgmtGroup.Use(AuthMiddleware(authService, []string{"products_admin"}))
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

func AuthMiddleware(authService services.AuthService, requiredDomains []string) func(context *gin.Context) {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		err, isValid := authService.Validate(token, requiredDomains)
		if err != nil {
			helpers.UnauthorizedRequest(err, context)
		} else if !isValid {
			helpers.UnauthorizedRequest(errors.New("not authorised"), context)
		}
	}
}
