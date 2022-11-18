package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/auth_service/internal/repositories"
	"github.com/manoamaro/microservices-store/auth_service/models"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"net/http"
)

type AuthController struct {
	r              *gin.Engine
	authRepository repositories.AuthRepository
}

func NewAuthController(r *gin.Engine, repository repositories.AuthRepository) *AuthController {
	return &AuthController{
		r:              r,
		authRepository: repository,
	}
}

func (a *AuthController) RegisterRoutes() {
	authRoute := a.r.Group("/auth")
	{
		authRoute.POST("/sign_up", a.signUpHandler)
		authRoute.POST("/sign_in", a.signInHandler)

		authorizedRoutes := authRoute.Group("/")
		authorizedRoutes.Use(func(c *gin.Context) {
			claims, token, err := a.authRepository.GetTokenFromRequest(c.Request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
			} else {
				c.Set("claims", claims)
				c.Set("token", token)
			}
		})
		authorizedRoutes.GET("/verify", a.verifyHandler)
		authorizedRoutes.DELETE("/invalidate", a.invalidateHandler)
	}
}

func (a *AuthController) signUpHandler(c *gin.Context) {
	request := &models.SignUpRequest{}
	if err := c.BindJSON(request); err != nil {
		helpers.BadRequest(err, c)
	} else if auth, err := a.authRepository.CreateAuth(request.Email, request.Password); err != nil {
		helpers.BadRequest(err, c)
	} else if signedString, err := a.authRepository.CreateToken(auth); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.Header("Authorization", "bearer "+signedString)
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func (a *AuthController) signInHandler(c *gin.Context) {
	request := &models.SignInRequest{}
	if err := c.BindJSON(request); err != nil {
		helpers.BadRequest(err, c)
	} else if auth, found := a.authRepository.Authenticate(request.Email, request.Password); !found {
		helpers.BadRequest(errors.New("auth not found"), c)
	} else if signedString, err := a.authRepository.CreateToken(auth); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.Header("Authorization", "bearer "+signedString)
		c.Status(http.StatusOK)
	}
}

func (a *AuthController) verifyHandler(c *gin.Context) {
	userClaims := c.MustGet("claims").(*models.UserClaims)
	c.JSON(http.StatusOK, gin.H{
		"audiences": userClaims.Audience,
		"flags":     userClaims.Flags,
	})
}

func (a *AuthController) invalidateHandler(c *gin.Context) {
	userClaims := c.MustGet("claims").(*models.UserClaims)
	token := c.MustGet("token").(string)
	if err := a.authRepository.InvalidateToken(userClaims, token); err != nil {
		helpers.UnauthorizedRequest(err, c)
	} else {
		c.Status(http.StatusOK)
	}
}
