package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"manoamaro.github.com/auth_service/internal/helpers"
	"manoamaro.github.com/auth_service/internal/repositories"
	"net/http"
	"time"

	"manoamaro.github.com/auth_service/internal"
)

var authRepository repositories.AuthRepository

func main() {

	authRepository = repositories.NewDefaultAuthRepository()

	r := gin.Default()
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/sign_up", signUpHandler)
		authRoute.POST("/sign_in", signInHandler)
		authorizedRoutes := authRoute.Group("/")
		authorizedRoutes.Use(func(c *gin.Context) {
			claims, token, err := authRepository.GetTokenFromRequest(c.Request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
			} else {
				c.Set("claims", claims)
				c.Set("token", token)
			}
		})
		authorizedRoutes.GET("/verify", verifyHandler)
		authorizedRoutes.DELETE("/invalidate", invalidateHandler)
	}

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func signUpHandler(c *gin.Context) {
	request := &internal.SignUpRequest{}
	if err := c.BindJSON(request); err != nil {
		handleError(err, c)
	} else if auth, err := authRepository.CreateAuth(request.Email, request.Password); err != nil {
		handleError(err, c)
	} else if signedString, err := authRepository.CreateToken(auth); err != nil {
		handleError(err, c)
	} else {
		c.Header("Authorization", "bearer "+signedString)
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func signInHandler(c *gin.Context) {
	request := &internal.SignInRequest{}
	if err := c.BindJSON(request); err != nil {
		handleError(err, c)
	} else if auth, found := authRepository.Authenticate(request.Email, request.Password); !found {
		handleError(errors.New("auth not found"), c)
	} else if signedString, err := authRepository.CreateToken(auth); err != nil {
		handleError(err, c)
	} else {
		c.Header("Authorization", "bearer "+signedString)
		c.Status(http.StatusOK)
	}
}

func verifyHandler(c *gin.Context) {
	userClaims := c.MustGet("claims").(*helpers.UserClaims)
	c.JSON(http.StatusOK, gin.H{
		"audiences": userClaims.Audience,
		"flags":     userClaims.Flags,
	})
}

func invalidateHandler(c *gin.Context) {
	userClaims := c.MustGet("claims").(*helpers.UserClaims)
	token := c.MustGet("token").(string)
	if err := authRepository.InvalidateToken(userClaims, token); err != nil {
		handleError(err, c)
	} else {
		c.Status(http.StatusOK)
	}
}

func handleError(err error, c *gin.Context) {
	log.Println(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status": err.Error(),
	})
}
