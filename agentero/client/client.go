package main

import (
	"fmt"
	"log"

	"github.com/agentero-go/policy_holder/policy_holder_pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	client := policy_holder_pb.NewPolicyHoldersServiceClient(conn)
	fmt.Printf("Created client: %f\n\n", client)

	g := gin.Default()

	if err := g.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
