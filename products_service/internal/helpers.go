package internal

import "log"

func FailOnError(err error) {
	if err != nil {
		log.Fatalf("%s")
	}
}
