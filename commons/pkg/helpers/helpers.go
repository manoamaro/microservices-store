package helpers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func GetEnv(name string, fallback string) (value string) {
	value, found := os.LookupEnv(name)
	if !found {
		log.Printf("Env %s not found. Fallback to %s", name, fallback)
		value = fallback
	}
	return
}

type ENV string

const (
	DEV  ENV = "dev"
	TEST ENV = "test"
	PROD ENV = "prod"
)

func IsEnvironment(env ENV) bool {
	current := GetEnv("environment", string(DEV))
	return ENV(current) == env
}

func GetFromContext[T any](c *gin.Context, key string) T {
	return c.MustGet(key).(T)
}

func BadRequest(err error, c *gin.Context) {
	handleError(err, c, http.StatusBadRequest)
}

func UnauthorizedRequest(err error, c *gin.Context) {
	handleError(err, c, http.StatusUnauthorized)
}

func handleError(err error, c *gin.Context, status int) {
	log.Println(err)
	c.AbortWithStatusJSON(status, gin.H{
		"status": err.Error(),
	})
}
