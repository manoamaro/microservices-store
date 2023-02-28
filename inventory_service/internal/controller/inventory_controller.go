package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/inventory_service/internal/repository"
)

type InventoryController struct {
	r                   *gin.Engine
	authService         infra.AuthService
	inventoryRepository repository.InventoryRepository
}

func NewInventoryController(r *gin.Engine, authService infra.AuthService, inventoryRepository repository.InventoryRepository) InventoryController {
	return InventoryController{
		r:                   r,
		authService:         authService,
		inventoryRepository: inventoryRepository,
	}
}
func (a *InventoryController) RegisterRoutes() {
	public := a.r.Group("/public")
	public.GET("/inventory/:product_id", a.amountOfHandler)

	internal := a.r.Group("/internal", helpers.AuthMiddleware(a.authService, "inventory"))
	internal.POST("/inventory/add", a.addHandler)
	internal.POST("/inventory/subtract", a.subtractHandler)
}

func (a *InventoryController) amountOfHandler(c *gin.Context) {
	var uri struct {
		ProductId string `uri:"product_id" binding:"required"`
	}
	if err := c.BindUri(&uri); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	amount := a.inventoryRepository.AmountOf(uri.ProductId)
	c.JSON(200, gin.H{"amount": amount})
}

func (a *InventoryController) addHandler(c *gin.Context) {
	var request struct {
		ProductId string `json:"product_id" binding:"required"`
		Amount    uint   `json:"amount" binding:"required"`
	}
	if err := c.BindJSON(&request); err == nil {
		amount := a.inventoryRepository.Add(request.ProductId, request.Amount)
		c.JSON(200, gin.H{"amount": amount})
	}
}

func (a *InventoryController) subtractHandler(c *gin.Context) {
	var request struct {
		ProductId string `json:"product_id" binding:"required"`
		Amount    uint   `json:"amount" binding:"required"`
	}
	if err := c.BindJSON(&request); err == nil {
		amount := a.inventoryRepository.Subtract(request.ProductId, request.Amount)
		c.JSON(200, gin.H{"amount": amount})
	}
}
