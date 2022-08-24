package internal

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"manoamaro.github.com/commons/pkg"
	"manoamaro.github.com/products_service/internal/controller"
	"manoamaro.github.com/products_service/internal/repository"
	"manoamaro.github.com/products_service/internal/service"
)

const ProductsServiceDatabase = "ProductsService"

type Application struct {
	ProductsRepository repository.ProductsRepository
	AuthService        service.AuthService
}

func NewApplication() *Application {
	authUrl := pkg.GetEnv("AUTH_URL", "http://localhost:8081/auth")
	mongoUrl := pkg.GetEnv("MONGO_URL", "mongodb://test:test@localhost:27017/?maxPoolSize=20&w=majority")
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		panic(err)
	}
	mongoDB := mongoClient.Database(ProductsServiceDatabase)
	return &Application{
		ProductsRepository: repository.NewDefaultProductsRepository(mongoDB),
		AuthService:        service.NewDefaultAuthService(authUrl),
	}
}

func (a *Application) SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("productsRepository", a.ProductsRepository)
		c.Set("authService", a.AuthService)
	})

	controller.ProductController(r)
	controller.AdminProductController(r)

	return r
}
