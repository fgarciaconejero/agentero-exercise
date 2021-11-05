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

var users = []policy_holder_pb.PolicyHolder{
	{
		Name:         "user1",
		MobileNumber: "1234567890",
	},
	{
		Name:         "user2",
		MobileNumber: "123 456 7891",
	},
	{
		Name:         "user3",
		MobileNumber: "(123) 456 7892",
	},
	{
		Name:         "user4",
		MobileNumber: "(123) 456-7893",
	},
	{
		Name:         "user5",
		MobileNumber: "123-456-7894",
	}}
var policies = []policy_holder_pb.InsurancePolicy{
	{
		MobileNumber: "1234567890",
		Premium:      2000,
		Type:         "homeowner",
	},
	{
		MobileNumber: "123 456 7891",
		Premium:      200,
		Type:         "renter",
	},
	{
		MobileNumber: "123-456-7892",
		Premium:      1500,
		Type:         "homeowner",
	},
	{
		MobileNumber: "(123) 456-7893",
		Premium:      155,
		Type:         "personal_auto",
	},
	{
		MobileNumber: "123-456-7894",
		Premium:      1000,
		Type:         "homeowner",
	},
	{
		MobileNumber: "123-456-7890",
		Premium:      500,
		Type:         "personal_auto",
	},
	{
		MobileNumber: "1234567892",
		Premium:      100,
		Type:         "life",
	},
	{
		MobileNumber: "(123)456-7892",
		Premium:      200,
		Type:         "homeowner",
	},
}
