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

func MapTo[I interface{}, O interface{}](i []I, f func(I) O) []O {
	var output []O
	for _, el := range i {
		output = append(output, f(el))
	}
	return output
}
