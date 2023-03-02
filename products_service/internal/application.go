package internal

import (
	"context"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"

	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/products_service/internal/controller"
	"github.com/manoamaro/microservices-store/products_service/internal/mock"
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
	if helpers.IsEnvironment(helpers.DEV) {
		return newDevApplication()
	} else {
		return newProdApplication()
	}
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
		AuthService:        infra.NewDefaultAuthService(authUrl),
	}
}

func newDevApplication() *Application {
	return &Application{
		ProductsRepository: mock.NewProductsRepository(),
		AuthService:        &mock.AuthService{},
	}
}

func (a *Application) SetupRoutes() *gin.Engine {
	r := gin.Default()
	controller.NewProductController(r, a.AuthService, a.ProductsRepository)
	controller.NewAdminProductController(r, a.ProductsRepository, a.AuthService)
	return r
}
