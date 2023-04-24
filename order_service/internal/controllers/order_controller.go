package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/order_service/internal/repositories"
)

type OrderController struct {
	engine          *gin.Engine
	authService     infra.AuthService
	orderRepository repositories.OrderRepository
}

func NewOrderController(
	engine *gin.Engine,
	authService infra.AuthService,
	repository repositories.OrderRepository,
) infra.Controller {
	return &OrderController{
		engine:          engine,
		authService:     authService,
		orderRepository: repository,
	}
}

func (a *OrderController) RegisterRoutes() {
	public := a.engine.Group("/public")
	public.GET("/orders", a.listHandler)
	public.GET("/orders/:order_id", a.getHandler)

	internal := a.engine.Group("/internal", helpers.AuthMiddleware(a.authService, "order"))
	internal.GET("/orders", a.listInternalHandler)
}

func (a *OrderController) listHandler(c *gin.Context) {

}

func (a *OrderController) getHandler(c *gin.Context) {

}

// handle a request to list all orders from db
func (a *OrderController) listInternalHandler(c *gin.Context) {

}
