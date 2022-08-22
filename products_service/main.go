package main

import (
	"log"
	"manoamaro.github.com/products_service/internal"
	"manoamaro.github.com/products_service/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var productsRepository models.ProductsRepository
var authService internal.AuthService

func main() {

	productsRepository = models.NewProductsRepository()
	authService = internal.NewAuthService()

	r := gin.Default()

	productsGroup := r.Group("/products")

	productsGroup.GET("/", ListProductsHandler)

	mgmtGroup := productsGroup.Group("/mgmt")
	mgmtGroup.Use(func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		err, isValid := authService.Validate(token)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
		} else if !isValid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "not authorized"})
		}
	})
	mgmtGroup.GET("/list", ListProductsHandler)
	mgmtGroup.POST("/create", PostProductsHandler)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handleError(err error, c *gin.Context) {
	log.Println(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status": err.Error(),
	})
}

func ListProductsHandler(c *gin.Context) {
	if products, err := productsRepository.ListProducts(); err != nil {
		handleError(err, c)
	} else {
		c.JSON(http.StatusOK, products)
	}
}

func PostProductsHandler(c *gin.Context) {
	newProduct := models.Product{}
	if err := c.BindJSON(&newProduct); err != nil {
		handleError(err, c)
	} else if savedProduct, err := productsRepository.InsertProduct(newProduct); err != nil {
		handleError(err, c)
	} else {
		c.JSON(http.StatusOK, savedProduct)
	}
}
