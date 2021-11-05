package main

import (
	"log"
	"net/http"

	"github.com/agentero-exercise/agentero/resources/mocks"
	"github.com/gin-gonic/gin"
)

func main() {
	handleRequests()
}

func handleRequests() {
	g := gin.Default()

	g.GET("/users/:agentid", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, &mocks.Users)
	})

	g.GET("/policies/:agentId", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"policy_holder": mocks.Policies,
		})
	})

	if err := g.Run("localhost:8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
