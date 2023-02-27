package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/inventory_service/internal/repository"
)

type InventoryController struct {
	r                   *gin.Engine
	inventoryRepository repository.InventoryRepository
}

func NewInventoryController(r *gin.Engine, inventoryRepository repository.InventoryRepository) InventoryController {
	return InventoryController{
		r:                   r,
		inventoryRepository: inventoryRepository,
	}
}
func (a *InventoryController) RegisterRoutes() {
	a.r.GET("/inventory/:product_id", a.amountOfHandler)
	a.r.POST("/inventory/add", a.addHandler)
	a.r.POST("/inventory/subtract", a.subtractHandler)
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
	c.BindJSON(&request)
	amount := a.inventoryRepository.Add(request.ProductId, request.Amount)
	c.JSON(200, gin.H{"amount": amount})
}

func (a *InventoryController) subtractHandler(c *gin.Context) {
	var request struct {
		ProductId string `json:"product_id" binding:"required"`
		Amount    uint   `json:"amount" binding:"required"`
	}
	c.BindJSON(&request)
	amount := a.inventoryRepository.Subtract(request.ProductId, request.Amount)
	c.JSON(200, gin.H{"amount": amount})
}
