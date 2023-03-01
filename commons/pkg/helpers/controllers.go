package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
)

func AuthMiddleware(authService infra.AuthService, requiredDomains ...string) func(context *gin.Context) {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		isValid, err := authService.Validate(token, requiredDomains...)
		if err != nil {
			UnauthorizedRequest(err, context)
		} else if !isValid {
			UnauthorizedRequest(errors.New("not authorised"), context)
		}
	}
}
