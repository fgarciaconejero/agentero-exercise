package main

import (
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

	// This for range is not done in the logical 'for i, v :=...' way because the mocks contain a Mux Lock
	// that not only throws a warning, but messes with the code logic
	for i := range mocks.Policies {
		if mocks.Policies[i].AgentId == agentId {
			mocks.Policies[i].MobileNumber = reg.ReplaceAllString(mocks.Policies[i].MobileNumber, "")
			ips = append(ips, &mocks.Policies[i])
		}
	}
	// fmt.Printf("Agent id: %v has ips: %vn", agentId, ips)
	return
}

func returnUsersByAgentId(agentId string) (phs []*protos.PolicyHolder) {
	reg := compileRegexp()

	// This for range is not done in the logical 'for i, v :=...' way because the mocks contain a Mux Lock
	// that not only throws a warning, but messes with the code logic
	for i := range mocks.Users {
		mocks.Users[i].MobileNumber = reg.ReplaceAllString(mocks.Users[i].MobileNumber, "")
		for x := range mocks.Policies {
			mocks.Policies[x].MobileNumber = reg.ReplaceAllString(mocks.Policies[x].MobileNumber, "")
			if (mocks.Policies[x].MobileNumber == mocks.Users[i].MobileNumber) && (mocks.Policies[x].AgentId == agentId) {
				phs = append(phs, &mocks.Users[i])
				break
			}
		}
	}
	return
}
