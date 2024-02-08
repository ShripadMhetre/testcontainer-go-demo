package utils

import (
	"github.com/google/uuid"
	"log"
)

func GenerateUUID() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal("FAILED TO CREATE NEW UUID", err)
	}

	return newUUID.String()
}
