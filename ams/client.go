package main

import (
	"log"
	"net/http"

	"github.com/agentero-go/policy_holder/policy_holder_pb"
	"github.com/gin-gonic/gin"
)

func main() {
	handleRequests()
}

func handleRequests() {
	g := gin.Default()

	g.GET("/users/:agentid", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, &users)
	})

	g.GET("/policies/:agentId", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"policy_holder": policies,
		})
	})

	if err := g.Run("localhost:8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

var users = []policy_holder_pb.PolicyHolder{}
var policies = []policy_holder_pb.InsurancePolicy{}
