package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/inventory_service/internal/repository"
)

type InventoryController struct {
	engine              *gin.Engine
	authService         infra.AuthService
	inventoryRepository repository.InventoryRepository
}

func NewInventoryController(
	engine *gin.Engine,
	authService infra.AuthService,
	inventoryRepository repository.InventoryRepository,
) infra.Controller {
	return &InventoryController{
		engine:              engine,
		authService:         authService,
		inventoryRepository: inventoryRepository,
	}
}

func (a *InventoryController) RegisterRoutes() {
	public := a.engine.Group("/public")
	public.GET("/inventory/:product_id", a.amountOfHandler)

	internal := a.engine.Group("/internal", helpers.AuthMiddleware(a.authService, "inventory"))
	internal.POST("/inventory/add", a.addHandler)
	internal.POST("/inventory/subtract", a.subtractHandler)
	internal.POST("/inventory/reserve", a.reserveHandler)
	internal.POST("/inventory/set", a.setHandler)
}

func (a *InventoryController) amountOfHandler(c *gin.Context) {
	var uri struct {
		ProductId string `uri:"product_id" binding:"required"`
	}
	if err := c.BindUri(&uri); err != nil {
		helpers.BadRequest(err, c)
	} else if amount, err := a.inventoryRepository.AmountOf(uri.ProductId); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(200, gin.H{"amount": amount})
	}
}

func (a *InventoryController) addHandler(c *gin.Context) {
	var request struct {
		ProductId string `json:"product_id" binding:"required"`
		Amount    uint   `json:"amount" binding:"required"`
	}
	if err := c.BindJSON(&request); err != nil {
		helpers.BadRequest(err, c)
	} else if amount, err := a.inventoryRepository.Add(request.ProductId, request.Amount); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(200, gin.H{"amount": amount})
	}
}

func (a *InventoryController) subtractHandler(c *gin.Context) {
	var request struct {
		ProductId string `json:"product_id" binding:"required"`
		Amount    uint   `json:"amount" binding:"required"`
	}
	if err := c.BindJSON(&request); err != nil {
		helpers.BadRequest(err, c)
	} else if amount, err := a.inventoryRepository.Subtract(request.ProductId, request.Amount); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(200, gin.H{"amount": amount})
	}
}

func (a *InventoryController) setHandler(c *gin.Context) {
	var request struct {
		ProductId string `json:"product_id" binding:"required"`
		Amount    uint   `json:"amount" binding:"required"`
	}
	if err := c.BindJSON(&request); err != nil {
		helpers.BadRequest(err, c)
	} else if amount, err := a.inventoryRepository.Set(request.ProductId, request.Amount); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(200, gin.H{"amount": amount})
	}
}

func (a *InventoryController) reserveHandler(c *gin.Context) {
	var request struct {
		ProductId string `json:"product_id" binding:"required"`
		CartId    string `json:"cart_id" binding:"required"`
		Amount    uint   `json:"amount" binding:"required"`
	}
	if err := c.BindJSON(&request); err != nil {
		helpers.BadRequest(err, c)
	} else if amount, err := a.inventoryRepository.Reserve(request.CartId, request.ProductId, request.Amount); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(200, gin.H{"amount": amount})
	}
}
