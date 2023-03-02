package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/collections"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
)

const (
	UserId        = "userId"
	UserAudiences = "userAudiences"
	UserFlags     = "userFlags"
)

func AuthMiddleware(authService infra.AuthService, requiredDomains ...string) func(context *gin.Context) {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		res, err := authService.Validate(token, requiredDomains...)
		if err != nil {
			UnauthorizedRequest(err, context)
		} else if len(requiredDomains) > 0 && !collections.ContainsAny(requiredDomains, res.Audiences) {
			UnauthorizedRequest(errors.New("not authorised"), context)
		} else {
			context.Set(UserId, res.UserId)
			context.Set(UserAudiences, res.Audiences)
			context.Set(UserFlags, res.Flags)
		}
	}
}
