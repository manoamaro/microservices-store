package internal

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/products_service/internal/controller"
	"github.com/manoamaro/microservices-store/products_service/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ProductsServiceDatabase = "ProductsService"

type Application struct {
	engine             *gin.Engine
	ProductsRepository repository.ProductsRepository
	AuthService        infra.AuthService
	controllers        []infra.Controller
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

	engine := gin.Default()
	// Enable CORS
	engine.Use(cors.New(infra.CorsConfig()))

	authService := infra.NewHttpAuthService(authUrl)
	productsRepository := repository.NewMongoDBProductsRepository(mongoDB)

	return &Application{
		engine:             engine,
		ProductsRepository: productsRepository,
		AuthService:        authService,
		controllers: []infra.Controller{
			controller.NewProductController(engine, authService, productsRepository),
			controller.NewAdminProductController(engine, productsRepository, authService),
		},
	}
}

func (a *Application) RegisterControllers() {
	for _, _controller := range a.controllers {
		_controller.RegisterRoutes()
	}
}

func (a *Application) Run(c chan error) {
	a.RegisterControllers()

	port := helpers.GetEnv("PORT", "8080")

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			c <- err
		}
	}()
}
