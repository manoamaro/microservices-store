package internal

import (
	"log"
	"os"
)

const ProductsServiceDatabase = "ProductsService"

func GetEnv(name string, fallback string) (value string) {
	value, found := os.LookupEnv(name)
	if !found {
		log.Printf("Env %s not found. Fallback to %s", name, fallback)
		value = fallback
	}
	return
}
