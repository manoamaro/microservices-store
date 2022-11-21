package main

import (
	"context"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo/mongodriver"
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservice-store/monolith/models"
	"github.com/manoamaro/microservice-store/monolith/repositories"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoUrl := helpers.GetEnv("MONGO_URL", "mongodb://admin:admin@localhost:27017/?maxPoolSize=20&w=majority")
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		panic(err)
	}
	mongoDB := mongoClient.Database("StoreMonolith")

	productsRepository := repositories.NewProductsRepository(mongoDB)
	usersRepository := repositories.NewUsersRepository(mongoDB)
	cartsRepository := repositories.NewCartsRepository(mongoDB)

	router := gin.Default()

	sessionsCollection := mongoDB.Collection("sessions")
	sessionStore := mongodriver.NewStore(sessionsCollection, 3600, true, []byte("secret"))
	router.Use(sessions.Sessions("session", sessionStore))

	authMiddleware := func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get("user")
		if userId == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Set("userId", userId)
		}
	}

	router.GET("/products", func(ctx *gin.Context) {
		var products []models.Product
		var err error

		query := ctx.Query("query")

		if len(query) > 3 {
			products, err = productsRepository.SearchProducts(query)
		} else {
			products, err = productsRepository.ListProducts()
		}

		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, products)
		}
	})

	router.GET("/products/:id", func(ctx *gin.Context) {
		objectId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else {
			product, err := productsRepository.GetProduct(objectId)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
			} else {
				ctx.JSON(200, product)
			}
		}
	})

	type CreateReviewRequest struct {
		Stars   int    `form:"stars" json:"stars" xml:"stars"  binding:"required"`
		Comment string `form:"comment" json:"comment" xml:"comment" binding:"required"`
	}

	router.POST("/products/:id/review", authMiddleware, func(ctx *gin.Context) {
		var request CreateReviewRequest

		if err := ctx.BindJSON(&request); err == nil {
			if review, err := productsRepository.CreateReview(ctx.Param("id"), ctx.GetString("userId"), request.Stars, request.Comment); err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
			} else {
				ctx.JSON(200, review)
			}
		}
	})

	router.GET("/cart", authMiddleware, func(ctx *gin.Context) {
		userId := ctx.GetString("userId")
		if cart, err := cartsRepository.GetCartForUser(userId); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else {
			ctx.JSON(200, cart)
		}
	})

	type PutCartRequest struct {
		ProductId string `form:"product_id" json:"product_id" xml:"product_id"  binding:"required"`
		Quantity  int    `form:"quantity" json:"quantity" xml:"quantity" binding:"required"`
	}
	router.PUT("/cart", authMiddleware, func(ctx *gin.Context) {
		request := &PutCartRequest{}

		if err := ctx.BindJSON(request); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else if productId, err := primitive.ObjectIDFromHex(request.ProductId); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else if product, err := productsRepository.GetProduct(productId); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else if cart, err := cartsRepository.AddProduct(ctx.GetString("userId"), *product, request.Quantity); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else {
			ctx.JSON(200, cart)
		}
	})

	router.PUT("/cart/:itemId", authMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.DELETE("/cart/:itemId", authMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.PUT("/cart/address", authMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/cart/place", authMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/orders", authMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/orders/:id", authMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.PUT("/orders/:id", authMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	type Login struct {
		Email    string `form:"email" json:"email" xml:"email"  binding:"required"`
		Password string `form:"password" json:"password" xml:"password" binding:"required"`
	}

	router.POST("/user/login", func(ctx *gin.Context) {
		var json Login

		if err := ctx.BindJSON(&json); err == nil {
			if user, err := usersRepository.Login(json.Email, json.Password); err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
			} else {
				session := sessions.Default(ctx)
				session.Set("user", user.Id.Hex())
				session.Save()
				ctx.JSON(http.StatusOK, user)
			}
		}
	})

	router.POST("/user/signup", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run()
}
