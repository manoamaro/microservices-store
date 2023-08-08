package helpers

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/manoamaro/microservices-store/commons/pkg/collections"
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
)

const (
	UserId        = "userId"
	UserAudiences = "userAudiences"
	UserFlags     = "userFlags"
)

var ErrNotAuthorised = errors.New("not authorised")

func AuthMiddleware(authService infra.AuthService, requiredDomains ...string) func(context *gin.Context) {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		res, err := authService.Validate(token, requiredDomains...)
		if err != nil {
			UnauthorizedRequest(err, context)
		} else if len(requiredDomains) > 0 && !collections.ContainsAny(requiredDomains, res.Audiences) {
			UnauthorizedRequest(ErrNotAuthorised, context)
		} else {
			context.Set(UserId, res.UserId)
			context.Set(UserAudiences, res.Audiences)
			context.Set(UserFlags, res.Flags)
		}
	}
}

func GetHost(ctx *gin.Context) string {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, ctx.Request.Host)
}

func CorsConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, "http://localhost:3000")
	return corsConfig
}
