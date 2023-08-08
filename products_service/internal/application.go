package internal

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"

	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/products_service/internal/controller"
	"github.com/manoamaro/microservices-store/products_service/internal/repository"
	"github.com/manoamaro/microservices-store/products_service/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ProductsServiceDatabase = "ProductsService"

type Application struct {
	ProductsRepository repository.ProductsRepository
	AuthService        infra.AuthService
	InventoryService   services.InventoryService
}

func NewApplication() *Application {
	return newProdApplication()
}

func newProdApplication() *Application {
	authUrl := helpers.GetEnv("AUTH_SERVICE_URL", "http://localhost:8081")
	mongoUrl := helpers.GetEnv("MONGO_URL", "mongodb://test:test@localhost:27017/?maxPoolSize=20&w=majority")
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		panic(err)
	}
	mongoDB := mongoClient.Database(ProductsServiceDatabase)
	return &Application{
		ProductsRepository: repository.NewDefaultProductsRepository(mongoDB),
		AuthService:        infra.NewHttpAuthService(authUrl),
	}
}

func (a *Application) SetupRoutes() *gin.Engine {
	engine := gin.Default()

	// Enable CORS
	engine.Use(cors.New(infra.CorsConfig()))

	controller.NewProductController(engine, a.AuthService, a.ProductsRepository)
	controller.NewAdminProductController(engine, a.ProductsRepository, a.AuthService)
	return engine
}
