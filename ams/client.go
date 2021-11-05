package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	handleRequests()
}

func handleRequests() {
	g := gin.Default()

	if err := g.Run("localhost:8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
