package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

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

func compileRegexp() (reg *regexp.Regexp) {
	// This regexp filters everything but numbers out
	reg, _ = regexp.Compile("[^0-9]+")
	return
}

func returnPoliciesByAgentId(agentId string) (ips []*protos.InsurancePolicy) {
	reg := compileRegexp()

	// This for range is not done in the logical 'for i, v :=...' way because lint throws a warning
	for _, v := range mocks.Policies {
		if v.AgentId == agentId {
			v.MobileNumber = reg.ReplaceAllString(v.MobileNumber, "")
			ips = append(ips, &v)
		}
	}
	return
}

func returnUsersByAgentId(agentId string) (phs []*protos.PolicyHolder) {
	reg := compileRegexp()
	fmt.Println("Here?")

	// This for range is not done in the logical 'for i, v :=...' way because lint throws a warning
	for _, v := range mocks.Users {
		for _, x := range mocks.Policies {
			if (x.MobileNumber == v.MobileNumber) && (x.AgentId == agentId) {
				v.MobileNumber = reg.ReplaceAllString(v.MobileNumber, "")
				phs = append(phs, &v)
				break
			}
		}
	}
	fmt.Println("Here!")
	return
}
