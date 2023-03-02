package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/products_service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type ProductController struct {
	authService        infra.AuthService
	productsRepository repository.ProductsRepository
}

func NewProductController(r *gin.Engine, authService infra.AuthService, productsRepository repository.ProductsRepository) *ProductController {
	controller := &ProductController{
		authService,
		productsRepository,
	}
	publicGroup := r.Group("/public")
	{
		publicGroup.GET("/", controller.listProductsHandler)
		publicGroup.GET("/:id", controller.getProductHandler)
		publicGroup.POST("/:id/review", helpers.AuthMiddleware(controller.authService), controller.postProductReviewHandler)
	}

	return controller
}

func (c *ProductController) listProductsHandler(ctx *gin.Context) {
	if products, err := c.productsRepository.ListProducts(); err != nil {
		helpers.BadRequest(err, ctx)
	} else {
		ctx.JSON(http.StatusOK, products)
	}
}

func (c *ProductController) getProductHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helpers.BadRequest(err, ctx)
		return
	}
	if product, err := c.productsRepository.GetProduct(objectID); err != nil {
		helpers.BadRequest(err, ctx)
	} else {
		ctx.JSON(http.StatusOK, product)
	}
}

type PostProductReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

func (c *ProductController) postProductReviewHandler(ctx *gin.Context) {
	productId := ctx.Param("id")
	userId := ctx.GetString(helpers.UserId)
	var req PostProductReviewRequest
	if err := ctx.BindJSON(&req); err != nil {
		helpers.BadRequest(err, ctx)
	} else if review, err := c.productsRepository.CreateReview(productId, userId, req.Rating, req.Comment); err != nil {
		helpers.BadRequest(err, ctx)
	} else {
		ctx.JSON(http.StatusOK, review)
	}
}
