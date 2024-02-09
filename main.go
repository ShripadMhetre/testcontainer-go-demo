package main

import (
	"log"
	"testcontainers-demo/app"
)

func main() {
	router := app.SetupRouter()

	log.Fatal(router.Run(":8080"))
}
