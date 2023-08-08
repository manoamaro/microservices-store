package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/manoamaro/microservices-store/auth_service/internal/use_cases"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
	"net/http"
)

type AuthController struct {
	engine              *gin.Engine
	signInUseCase       use_cases.SignInUseCase
	signUpUseCase       use_cases.SignUpUseCase
	verifyUseCase       use_cases.VerifyUseCase
	refreshTokenUseCase use_cases.RefreshTokenUseCase
}

func NewAuthController(
	r *gin.Engine,
	signInUseCase use_cases.SignInUseCase,
	signUpUseCase use_cases.SignUpUseCase,
	verifyUseCase use_cases.VerifyUseCase,
	refreshTokenUseCase use_cases.RefreshTokenUseCase,
) infra.Controller {
	controller := &AuthController{
		engine:              r,
		signInUseCase:       signInUseCase,
		signUpUseCase:       signUpUseCase,
		verifyUseCase:       verifyUseCase,
		refreshTokenUseCase: refreshTokenUseCase,
	}
	return controller
}

func (a *AuthController) RegisterRoutes() {
	public := a.engine.Group("/public")
	{
		public.POST("/sign_up", a.signUpHandler)
		public.POST("/sign_in", a.signInHandler)
		public.GET("/verify", a.verifyHandler)
		public.POST("/refresh", a.refreshTokenHandler)
	}
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUpResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (a *AuthController) signUpHandler(c *gin.Context) {
	req := &SignUpRequest{}
	if err := c.BindJSON(req); err != nil {
		helpers.BadRequest(err, c)
	} else if result, err := a.signUpUseCase.SignUp(use_cases.SignUpDTO{
		Email:         req.Email,
		PlainPassword: req.Password,
	}); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, SignInUpResponse{
			Token:        result.Token,
			RefreshToken: result.RefreshToken,
		})
	}
}

func (a *AuthController) signInHandler(c *gin.Context) {
	r := &SignInRequest{}
	if err := c.BindJSON(r); err != nil {
		helpers.BadRequest(err, c)
	} else if result, err := a.signInUseCase.SignIn(use_cases.SignInDTO{
		Email:         r.Email,
		PlainPassword: r.Password,
	}); err != nil {
		helpers.BadRequest(err, c)
	} else {
		c.JSON(http.StatusOK, SignInUpResponse{
			Token:        result.Token,
			RefreshToken: result.RefreshToken,
		})
	}
}

func (a *AuthController) verifyHandler(c *gin.Context) {
	if token, err := request.AuthorizationHeaderExtractor.ExtractToken(c.Request); err != nil {
		helpers.UnauthorizedRequest(err, c)
	} else if result, err := a.verifyUseCase.Verify(use_cases.VerifyDTO{Token: token}); err != nil {
		helpers.UnauthorizedRequest(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"user_id":   result.ID,
			"audiences": result.Audience,
			"flags":     result.Flags,
		})
	}
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (a *AuthController) refreshTokenHandler(c *gin.Context) {
	req := RefreshTokenRequest{}
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(err, c)
	} else if res, err := a.refreshTokenUseCase.RefreshToken(use_cases.RefreshTokenDTO{RefreshToken: req.RefreshToken}); err != nil {
		helpers.UnauthorizedRequest(err, c)
	} else {
		c.JSON(http.StatusOK, SignInUpResponse{
			Token:        res.Token,
			RefreshToken: res.RefreshToken,
		})
	}
}
