package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"github.com/manoamaro/microservices-store/inventory_service/internal/use_cases/inventory"
)

type InventoryController struct {
	engine          *gin.Engine
	authService     infra.AuthService
	getUseCase      inventory.GetUseCase
	addUseCase      inventory.AddUseCase
	subtractUseCase inventory.SubtractUseCase
	reserveUseCase  inventory.ReserveUseCase
}

func NewInventoryController(
	engine *gin.Engine,
	authService infra.AuthService,
	getUseCase inventory.GetUseCase,
	addUseCase inventory.AddUseCase,
	subtractUseCase inventory.SubtractUseCase,
	reserveUseCase inventory.ReserveUseCase,
) infra.Controller {
	return &InventoryController{
		engine:          engine,
		authService:     authService,
		getUseCase:      getUseCase,
		addUseCase:      addUseCase,
		subtractUseCase: subtractUseCase,
		reserveUseCase:  reserveUseCase,
	}
}

func (a *InventoryController) RegisterRoutes() {
	public := a.engine.Group("/public")
	public.GET("/inventory/:product_id", a.amountOfHandler)

	internal := a.engine.Group("/internal", helpers.AuthMiddleware(a.authService, "inventory"))
	internal.POST("/inventory/add", a.addHandler)
	internal.POST("/inventory/subtract", a.subtractHandler)
	internal.POST("/inventory/reserve", a.reserveHandler)
}

func (a *InventoryController) amountOfHandler(c *gin.Context) {
	var uri struct {
		ProductId string `uri:"product_id" binding:"required"`
	}
	if err := c.BindUri(&uri); err != nil {
		helpers.BadRequest(err, c)
	} else if amount, err := a.getUseCase.Get(uri.ProductId); err != nil {
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
	} else if amount, err := a.addUseCase.Add(request.ProductId, request.Amount); err != nil {
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
	} else if amount, err := a.subtractUseCase.Subtract(request.ProductId, request.Amount); err != nil {
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
	} else if amount, err := a.reserveUseCase.Reserve(request.CartId, request.ProductId, request.Amount); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(200, gin.H{"amount": amount})
	}
}
