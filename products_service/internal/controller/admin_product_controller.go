package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/products_service/internal/models"
	"github.com/manoamaro/microservices-store/products_service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type AdminProductController struct {
	authService        infra.AuthService
	productsRepository repository.ProductsRepository
}

func NewAdminProductController(r *gin.Engine, productsRepository repository.ProductsRepository, authService infra.AuthService) *AdminProductController {
	controller := &AdminProductController{
		authService,
		productsRepository,
	}

	adminGroup := r.Group("/admin")
	{
		adminGroup.Use(helpers.AuthMiddleware(authService, "products_admin"))
		adminGroup.GET("/", controller.getProductsHandler)
		adminGroup.GET("/:id", controller.getProductHandler)
		adminGroup.POST("/", controller.postProductsHandler)
		adminGroup.POST("/:id/upload", controller.postProductImageHandler)
		adminGroup.PUT("/:id", controller.putProductsHandler)
		adminGroup.DELETE("/:id", controller.deleteProductsHandler)
	}
	return controller
}

func (c *AdminProductController) getProductsHandler(ctx *gin.Context) {
	if products, err := c.productsRepository.ListProducts(); err != nil {
		helpers.BadRequest(err, ctx)
	} else {
		ctx.JSON(http.StatusOK, products)
	}
}

func (c *AdminProductController) getProductHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	if objectID, err := primitive.ObjectIDFromHex(id); err != nil {
		helpers.BadRequest(err, ctx)
	} else if product, err := c.productsRepository.GetProduct(objectID); err != nil {
		helpers.BadRequest(err, ctx)
	} else {
		ctx.JSON(http.StatusOK, product)
	}
}

func (c *AdminProductController) postProductsHandler(ctx *gin.Context) {
	newProduct := models.Product{}
	if err := ctx.BindJSON(&newProduct); err != nil {
		helpers.BadRequest(err, ctx)
	} else if savedProduct, err := c.productsRepository.InsertProduct(newProduct); err != nil {
		helpers.BadRequest(err, ctx)
	} else {
		ctx.JSON(http.StatusCreated, savedProduct)
	}
}

func (c *AdminProductController) putProductsHandler(ctx *gin.Context) {
	newProduct := models.Product{}
	if err := ctx.BindJSON(&newProduct); err != nil {
		helpers.BadRequest(err, ctx)
	} else if savedProduct, err := c.productsRepository.UpdateProduct(newProduct); err != nil {
		helpers.BadRequest(err, ctx)
	} else {
		ctx.JSON(http.StatusCreated, savedProduct)
	}
}

func (c *AdminProductController) postProductImageHandler(ctx *gin.Context) {

	id := ctx.Param("id")
	if objectID, err := primitive.ObjectIDFromHex(id); err != nil {
		helpers.BadRequest(err, ctx)
	} else if product, err := c.productsRepository.GetProduct(objectID); err != nil {
		helpers.BadRequest(err, ctx)
	} else {
		// Multipart form
		form, _ := ctx.MultipartForm()
		files := form.File["images[]"]

		for _, file := range files {
			log.Println(file.Filename)
			// Upload the file to specific dst.
			imagePath := fmt.Sprintf("uploaded/%s", uuid.New().String())
			if err := ctx.SaveUploadedFile(file, imagePath); err != nil {
				helpers.BadRequest(err, ctx)
			} else if _, err := c.productsRepository.AddImage(product.Id, imagePath); err != nil {
				helpers.BadRequest(err, ctx)
			}
		}
		ctx.Status(http.StatusCreated)
	}
}

func (c *AdminProductController) deleteProductsHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	if objectID, err := primitive.ObjectIDFromHex(id); err != nil {
		helpers.BadRequest(err, ctx)
	} else if _, err := c.productsRepository.DeleteProduct(objectID); err != nil {
		helpers.BadRequest(err, ctx)
	} else {
		ctx.Status(http.StatusOK)
	}
}
