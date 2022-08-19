package internal

import (
	"log"
	"os"
)

func FailOnError(err error) {
	if err != nil {
		log.Fatalf("%s")
	}
}

func GetENV(name string, fallback string) (value string) {
	value, found := os.LookupEnv(name)
	if !found {
		log.Printf("Env %s not found. Fallback to %s", name, fallback)
		value = fallback
	}
	return
}
