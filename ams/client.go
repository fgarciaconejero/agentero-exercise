package main

import (
	"log"
	"net/http"

	"github.com/agentero-exercise/ams/resources/mocks"
	"github.com/gin-gonic/gin"
)

func main() {
	handleRequests()
}

func handleRequests() {
	g := gin.Default()

	// TODO: Make it so that this endpoint and the one below filter out by agentId
	g.GET("/users/:agentid", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mocks.Users)
	})

	g.GET("/policies/:agentId", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mocks.Policies)
	})

	g.GET("/agents", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mocks.Agents)
	})

	if err := g.Run("localhost:8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
