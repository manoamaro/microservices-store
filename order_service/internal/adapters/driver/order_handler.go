package driver_adapters

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/order_service/internal/core/application"
	driver_ports "github.com/manoamaro/microservices-store/order_service/internal/core/ports/driver"
	"strconv"
)

type ginOrderHandler struct {
	engine       *gin.Engine
	orderService application.OrderService
}

func NewGinOrderHandler(engine *gin.Engine, orderService application.OrderService) driver_ports.OrderApi {
	return &ginOrderHandler{
		engine:       engine,
		orderService: orderService,
	}
}

func (handler *ginOrderHandler) GetOrderHandler(c context.Context) {
	ctx := c.Value("ginContext").(*gin.Context)
	if orderId, err := strconv.Atoi(ctx.Param("order_id")); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid order id"})
	} else if order, err := handler.orderService.GetOrder(uint(orderId)); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(200, order)
	}
}
