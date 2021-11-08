package service_test

import (
	"errors"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/agentero-exercise/agentero/domain/models"
	"github.com/agentero-exercise/agentero/resources/mocks"
	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/agentero-exercise/agentero/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// This function servers as a BeforeEach() / Init() / SetUp() to run before tests and initialize stuff
// In this case, we need the AmsMockApi to be initialized globally or there will be port conflicts
func TestMain(m *testing.M) {
	initializeAmsMockApi()
	code := m.Run()

	os.Exit(code)
}

var getPolicyHoldersFromAmsTestingParameters = []struct {
	name       string
	id         string
	repository mockRepository
	expected   []*protos.PolicyHolder
	err        error
}{
	{
		"successful",
		"some-agent-id",
		mockRepository{},
		mocks.Users,
		nil,
	},
	{
		"isGetUsersByIdError true",
		"amsReturnUsers error",
		mockRepository{},
		nil,
		errors.New("HTTP 400"),
	},
}

func TestGetPolicyHoldersFromAms(t *testing.T) {
	for _, tt := range getPolicyHoldersFromAmsTestingParameters {
		s := service.NewService(&tt.repository)
		res, err := s.GetPolicyHoldersFromAms(tt.id)
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
		assert.Equal(t, res, tt.expected)
	}
}

var getInsurancePoliciesFromAmsTestingParameters = []struct {
	name       string
	id         string
	repository mockRepository
	expected   []*protos.InsurancePolicy
	err        error
}{
	{
		"successful",
		"some-agent-id",
		mockRepository{},
		mocks.Policies,
		nil,
	},
	{
		"isGetPoliciesByIdError true",
		"amsReturnPolicies error",
		mockRepository{},
		nil,
		errors.New("HTTP 400"),
	},
}

func TestGetInsurancePoliciesFromAms(t *testing.T) {
	for _, tt := range getInsurancePoliciesFromAmsTestingParameters {
		s := service.NewService(&tt.repository)
		res, err := s.GetInsurancePoliciesFromAms(tt.id)
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
		assert.Equal(t, res, tt.expected)
	}
}

var getInsuranceAgentsFromAmsTestingParameters = []struct {
	name       string
	id         string
	repository mockRepository
	expected   []*models.Agent
	err        error
}{
	{
		"successful",
		"some-agent-id",
		mockRepository{},
		mocks.Agents,
		nil,
	},
}

func TestGetInsuranceAgentsFromAms(t *testing.T) {
	for _, tt := range getInsuranceAgentsFromAmsTestingParameters {
		s := service.NewService(&tt.repository)
		res, err := s.GetInsuranceAgentsFromAms()
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
		assert.Equal(t, res, tt.expected)
	}
}

// Helper function to mock the AMS Mock Api we'll need to make the calls from the service
// Has slight variations of actual implementation to help testing, but should work pretty much the same
func initializeAmsMockApi() {
	g := gin.Default()

	g.GET("/agents", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mocks.Agents)
	})

	g.GET("/users/:agentId", func(ctx *gin.Context) {
		res, err := amsReturnUsers(ctx.Param("agentId"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, nil)
		} else {
			ctx.JSON(http.StatusOK, res)
		}
	})

	g.GET("/policies/:agentId", func(ctx *gin.Context) {
		res, err := amsReturnPolicies(ctx.Param("agentId"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, nil)
		} else {
			ctx.JSON(http.StatusOK, res)
		}
	})

	go func() {
		if err := g.Run("localhost:8081"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
		time.Sleep(1 * time.Second)
	}()
}

// Helper function to give tests the possibility of failure from the AMS endpoints
func amsReturnUsers(isError string) ([]*protos.PolicyHolder, error) {
	if isError == "amsReturnUsers error" {
		return nil, errors.New("HTTP 400")
	}
	return mocks.Users, nil
}

// Helper function to give tests the possibility of failure from the AMS endpoints
func amsReturnPolicies(isError string) ([]*protos.InsurancePolicy, error) {
	if isError == "amsReturnPolicies error" {
		return nil, errors.New("HTTP 400")
	}
	return mocks.Policies, nil
}

type mockRepository struct{}

func (*mockRepository) GetById(id string) ([]*protos.PolicyHolder, error) {
	return nil, nil
}

func (*mockRepository) GetByMobileNumber(id string) (*protos.PolicyHolder, error) {
	return nil, nil
}

func (*mockRepository) UpsertPolicyHolder(phs *protos.PolicyHolder) error {
	return nil
}

func (*mockRepository) UpsertInsurancePolicy(phs *protos.InsurancePolicy, agentId string) error {
	return nil
}

func (r *mockRepository) GetAllInsuranceAgentsIds() ([]string, error) {
	return []string{}, nil
}

func (r *mockRepository) UpsertInsuranceAgent(agents *models.Agent) error {
	return nil
}
