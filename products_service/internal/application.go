package internal

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"manoamaro.github.com/products_service/internal/models"
	"manoamaro.github.com/products_service/internal/repository"
	"manoamaro.github.com/products_service/internal/service"
	"net/http"
)

type Application struct {
	ProductsRepository repository.ProductsRepository
	AuthService        service.AuthService
}

func NewApplication() *Application {
	return &Application{
		ProductsRepository: repository.NewDefaultProductsRepository(),
		AuthService:        service.NewDefaultAuthService(),
	}
}

func (a *Application) SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("productsRepository", a.ProductsRepository)
		c.Set("authService", a.AuthService)
	})

	publicGroup := r.Group("/public")
	{
		publicGroup.GET("/", ListProductsHandler)
	}
	mgmtGroup := r.Group("/mgmt")
	{
		mgmtGroup.Use(AuthMiddleware([]string{"products_admin"}))
		mgmtGroup.GET("/list", ListProductsHandler)
		mgmtGroup.POST("/create", PostProductsHandler)
	}
	return r
}

func AuthMiddleware(requiredDomains []string) func(context *gin.Context) {
	return func(context *gin.Context) {
		auth := authService(context)
		token := context.GetHeader("Authorization")
		err, isValid := auth.Validate(token, requiredDomains)
		if err != nil {
			UnauthorizedRequest(err, context)
		} else if !isValid {
			UnauthorizedRequest(errors.New("not authorised"), context)
		}
	}
}

func ListProductsHandler(c *gin.Context) {
	productsRepository := productsRepository(c)
	if products, err := productsRepository.ListProducts(); err != nil {
		BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, products)
	}
}

func PostProductsHandler(c *gin.Context) {
	productsRepository := productsRepository(c)
	newProduct := models.Product{}
	if err := c.BindJSON(&newProduct); err != nil {
		BadRequest(err, c)
	} else if savedProduct, err := productsRepository.InsertProduct(newProduct); err != nil {
		BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, savedProduct)
	}
}

func get[T any](c *gin.Context, key string) T {
	return c.MustGet(key).(T)
}

func productsRepository(c *gin.Context) repository.ProductsRepository {
	return get[repository.ProductsRepository](c, "productsRepository")
}

func authService(c *gin.Context) service.AuthService {
	return get[service.AuthService](c, "authService")
}

func BadRequest(err error, c *gin.Context) {
	handleError(err, c, http.StatusBadRequest)
}

func UnauthorizedRequest(err error, c *gin.Context) {
	handleError(err, c, http.StatusUnauthorized)
}

func handleError(err error, c *gin.Context, status int) {
	log.Println(err)
	c.AbortWithStatusJSON(status, gin.H{
		"status": err.Error(),
	})
}
