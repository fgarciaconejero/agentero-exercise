package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := connectToServer()
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	client := newClient(conn)
	fmt.Println("Created client successfuly!")

	handleRequests(client)
}

func connectToServer() (*grpc.ClientConn, error) {
	return grpc.Dial("localhost:50051", grpc.WithInsecure())
}

func newClient(conn *grpc.ClientConn) protos.PolicyHoldersServiceClient {
	return protos.NewPolicyHoldersServiceClient(conn)
}

func handleRequests(client protos.PolicyHoldersServiceClient) {
	g := gin.Default()
	g.GET("/getById/:id", func(ctx *gin.Context) {
		req := &protos.GetContactAndPoliciesByIdRequest{
			InsuranceAgentId: ctx.Param("id"),
		}

		res, err := client.GetContactAndPoliciesById(ctx, req)
		if err != nil {
			log.Fatalf("Something went wrong while trying to get contact and policies by id: %v\n", err)
		}
		ctx.JSON(http.StatusOK, res.PolicyHolders)
	})

	g.GET("getByMobileNumber/:id/:mn", func(ctx *gin.Context) {
		req := &protos.GetContactsAndPoliciesByMobileNumberRequest{
			MobileNumber: ctx.Param("mn"),
		}

		res, err := client.GetContactsAndPoliciesByMobileNumber(ctx, req)
		if err != nil {
			log.Fatalf("Something went wrong while trying to get contact and policies by mobile number: %v\n", err)
		}
		ctx.JSON(http.StatusOK, res.PolicyHolder)
	})

	if err := g.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to run server: %v\n", err)
	}
}
