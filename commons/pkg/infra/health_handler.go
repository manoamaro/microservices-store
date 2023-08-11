package infra

import "github.com/gin-gonic/gin"

func HealthHandler(check func() error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := check(); err != nil {
			ctx.JSON(200, gin.H{"status": "ok"})
		} else {
			ctx.JSON(500, gin.H{"status": "error", "error": err.Error()})
		}
	}
}
