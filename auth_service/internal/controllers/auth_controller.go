package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/auth_service/internal/repositories"
	"github.com/manoamaro/microservices-store/auth_service/models"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"net/http"
	"strconv"
)

type AuthController struct {
	r              *gin.Engine
	authRepository repositories.AuthRepository
}

func NewAuthController(r *gin.Engine, repository repositories.AuthRepository) *AuthController {
	controller := &AuthController{
		r:              r,
		authRepository: repository,
	}
	controller.RegisterRoutes()
	return controller
}

func (a *AuthController) RegisterRoutes() {
	public := a.r.Group("/public")
	{
		public.POST("/sign_up", a.signUpHandler)
		public.POST("/sign_in", a.signInHandler)

		authorizedRoutes := public.Group("/")
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
		public.POST("/refresh", a.refreshTokenHandler)
	}
}

type SignResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (a *AuthController) signUpHandler(c *gin.Context) {
	request := &models.SignUpRequest{}
	if err := c.BindJSON(request); err != nil {
		helpers.BadRequest(err, c)
	} else if auth, err := a.authRepository.CreateAuth(request.Email, request.Password); err != nil {
		helpers.BadRequest(err, c)
	} else if accessToken, refreshToken, err := a.authRepository.CreateTokens(auth.ID); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, SignResponse{
			Token:        accessToken,
			RefreshToken: refreshToken,
		})
	}
}

func (a *AuthController) signInHandler(c *gin.Context) {
	request := &models.SignInRequest{}
	if err := c.BindJSON(request); err != nil {
		helpers.BadRequest(err, c)
	} else if auth, found := a.authRepository.Authenticate(request.Email, request.Password); !found {
		helpers.BadRequest(errors.New("auth not found"), c)
	} else if accessToken, refreshToken, err := a.authRepository.CreateTokens(auth.ID); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, SignResponse{
			Token:        accessToken,
			RefreshToken: refreshToken,
		})
	}
}

func (a *AuthController) verifyHandler(c *gin.Context) {
	userClaims := c.MustGet("claims").(*models.UserClaims)
	token := c.MustGet("token").(string)
	if a.authRepository.IsInvalidatedToken(token) {
		helpers.UnauthorizedRequest(fmt.Errorf("token invalidated"), c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"user_id":   userClaims.ID,
			"audiences": userClaims.Audience,
			"flags":     userClaims.Flags,
		})
	}
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (a *AuthController) refreshTokenHandler(c *gin.Context) {
	request := RefreshTokenRequest{}
	if err := c.BindJSON(&request); err != nil {
		helpers.BadRequest(err, c)
	} else if claims, err := a.authRepository.GetClaimsFromRefreshToken(request.RefreshToken); err != nil {
		helpers.UnauthorizedRequest(err, c)
	} else if a.authRepository.IsInvalidatedToken(request.RefreshToken) {
		helpers.UnauthorizedRequest(fmt.Errorf("token invalidated"), c)
	} else if authId, err := strconv.ParseUint(claims.ID, 10, 32); err != nil {
		helpers.BadRequest(err, c)
	} else if accessToken, refreshToken, err := a.authRepository.CreateTokens(uint(authId)); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, SignResponse{
			Token:        accessToken,
			RefreshToken: refreshToken,
		})
	}
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
