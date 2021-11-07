package main

import (
	"log"
	"net/http"

	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/agentero-exercise/ams/resources/mocks"
	"github.com/gin-gonic/gin"
)

func main() {
	handleRequests()
}

func handleRequests() {
	g := gin.Default()

	g.GET("/users/:agentId", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, returnUsersByAgentId(ctx.Param("agentId")))
	})

	g.GET("/policies/:agentId", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, returnPoliciesByAgentId(ctx.Param("agentId")))
	})

	g.GET("/agents", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mocks.Agents)
	})

	if err := g.Run("localhost:8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func returnPoliciesByAgentId(agentId string) (ips []*protos.InsurancePolicy) {
	// This for range is not done in the logical 'for i, v :=...' way because lint throws a warning
	for i := range mocks.Policies {
		if mocks.Policies[i].AgentId == agentId {
			ips = append(ips, &mocks.Policies[i])
		}
	}
	return
}

func returnUsersByAgentId(agentId string) (phs []*protos.PolicyHolder) {
	// This for range is not done in the logical 'for i, v :=...' way because lint throws a warning
	for i := range mocks.Users {
		for x := range mocks.Policies {
			if (mocks.Policies[x].MobileNumber == mocks.Users[i].MobileNumber) && (mocks.Policies[x].AgentId == agentId) {
				phs = append(phs, &mocks.Users[i])
				break
			}
		}
	}
	return
}
