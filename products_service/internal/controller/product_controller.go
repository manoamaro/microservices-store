package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"manoamaro.github.com/commons/pkg"
	"manoamaro.github.com/products_service/internal/models"
	"manoamaro.github.com/products_service/internal/repository"
	"manoamaro.github.com/products_service/internal/service"
	"net/http"
)

func ProductController(r *gin.Engine) {
	publicGroup := r.Group("/public")
	{
		publicGroup.GET("/", listProductsHandler)
	}
}

func AdminProductController(r *gin.Engine) {
	mgmtGroup := r.Group("/admin")
	{
		mgmtGroup.Use(AuthMiddleware([]string{"products_admin"}))
		mgmtGroup.GET("/list", listProductsHandler)
		mgmtGroup.POST("/create", postProductsHandler)
	}
}

func listProductsHandler(c *gin.Context) {
	productsRepository := ProductsRepository(c)
	if products, err := productsRepository.ListProducts(); err != nil {
		pkg.BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, products)
	}
}

func postProductsHandler(c *gin.Context) {
	productsRepository := ProductsRepository(c)
	newProduct := models.Product{}
	if err := c.BindJSON(&newProduct); err != nil {
		pkg.BadRequest(err, c)
	} else if savedProduct, err := productsRepository.InsertProduct(newProduct); err != nil {
		pkg.BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, savedProduct)
	}
}

func AuthMiddleware(requiredDomains []string) func(context *gin.Context) {
	return func(context *gin.Context) {
		auth := AuthService(context)
		token := context.GetHeader("Authorization")
		err, isValid := auth.Validate(token, requiredDomains)
		if err != nil {
			pkg.UnauthorizedRequest(err, context)
		} else if !isValid {
			pkg.UnauthorizedRequest(errors.New("not authorised"), context)
		}
	}
}

func ProductsRepository(c *gin.Context) repository.ProductsRepository {
	return pkg.GetFromContext[repository.ProductsRepository](c, "productsRepository")
}

func AuthService(c *gin.Context) service.AuthService {
	return pkg.GetFromContext[service.AuthService](c, "authService")
}
