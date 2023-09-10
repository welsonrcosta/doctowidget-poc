package util

import (
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal("An error has occured:", err)
	}
}
