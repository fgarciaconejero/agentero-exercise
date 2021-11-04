package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/agentero-go/policy_holder/policy_holder_pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to server: %v\n", err)
	}
	defer conn.Close()

	client := policy_holder_pb.NewPolicyHoldersServiceClient(conn)
	fmt.Printf("Created client: %f\n", client)

	g := gin.Default()

	g.GET("/getById/:id", func(ctx *gin.Context) {
		req := &policy_holder_pb.GetContactAndPoliciesByIdRequest{
			InsuranceAgentId: ctx.Param("id"),
		}

		res, err := client.GetContactAndPoliciesById(ctx, req)
		if err != nil {
			log.Fatalf("Something went wrong while trying to get contact and policies by id: %v\n", err)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(res.PolicyHolders),
		})
	})

	if err := g.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to run server: %v\n", err)
	}
}
