package main

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/agentero-exercise/agentero/resources/mocks"
	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/gin-gonic/gin"
)

func TestGetContactAndPoliciesById(t *testing.T) {
	initializeAmsMockApi()

	s := &server{}
	req := protos.GetContactAndPoliciesByIdRequest{
		InsuranceAgentId: "some-id",
	}

	res, err := s.GetContactAndPoliciesById(context.Background(), &req)
	if err != nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}
}

func TestGetContactsAndPoliciesByMobileNumber(t *testing.T) {
	s := &server{}
	req := protos.GetContactsAndPoliciesByMobileNumberRequest{
		MobileNumber: "some-mobile-number",
	}
	res, err := s.GetContactsAndPoliciesByMobileNumber(context.Background(), &req)
	if res != nil || err != nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}
}

func initializeAmsMockApi() {
	g := gin.Default()

	g.GET("/users/:agentid", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, &mocks.Users)
	})

	g.GET("/policies/:agentId", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"policy_holder": mocks.Policies,
		})
	})

	go func() {
		if err := g.Run("localhost:8081"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
		time.Sleep(5 * time.Second)
	}()

}
