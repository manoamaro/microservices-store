package internal

import (
	"context"
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
	AuthService        services.AuthService
}

func NewApplication() *Application {
	if helpers.IsEnvironment(helpers.DEV) {
		return newDevApplication()
	} else {
		return newProdApplication()
	}
}

func newProdApplication() *Application {
	authUrl := helpers.GetEnv("AUTH_URL", "http://localhost:8081/auth")
	mongoUrl := helpers.GetEnv("MONGO_URL", "mongodb://test:test@localhost:27017/?maxPoolSize=20&w=majority")
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		panic(err)
	}
	mongoDB := mongoClient.Database(ProductsServiceDatabase)
	return &Application{
		ProductsRepository: repository.NewDefaultProductsRepository(mongoDB),
		AuthService:        services.NewDefaultAuthService(authUrl),
	}
}

func newDevApplication() *Application {
	return &Application{
		ProductsRepository: mock.NewMockProductsRepository(),
		AuthService:        &mock.MockAuthService{},
	}
}

func (a *Application) SetupRoutes() *gin.Engine {
	r := gin.Default()
	controller.NewProductController(r, a.ProductsRepository)
	controller.NewAdminProductController(r, a.ProductsRepository, a.AuthService)
	return r
}
