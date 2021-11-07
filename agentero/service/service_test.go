package service_test

import (
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/agentero-exercise/agentero/domain/models"
	"github.com/agentero-exercise/agentero/resources/mocks"
	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/agentero-exercise/agentero/service"
	"github.com/gin-gonic/gin"
)

var getPolicyHoldersFromAmsTestingParameters = []struct {
	name       string
	id         string
	repository mockRepository
	expected   []*protos.PolicyHolder
	err        error
}{}

func TestGetPolicyHoldersFromAms(t *testing.T) {
	initializeAmsMockApi()
	for _, tt := range getPolicyHoldersFromAmsTestingParameters {
		s := service.NewService(&tt.repository)
		res, err := s.GetPolicyHoldersFromAms(tt.id)
		// Lint warns not to use DeepEqual on error, but every other way doesn't work or panics because
		// in the case of the error being nil there is a nil pointer exception
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Test '%v' failed! err: %v, expected: %v\n", tt.name, err, tt.err)
		}
		if !reflect.DeepEqual(res, tt.expected) {
			t.Errorf("Test '%v' failed! \nres: %v,\n expected: %v\n", tt.name, res, tt.expected)
		}
	}
}

// func TestGetInsurancePoliciesFromAms(t *testing.T) {
// 	initializeAmsMockApi()

// 	s := service.NewService(&mockRepository{})

// 	res, err := s.GetInsurancePoliciesFromAms("some-agent-id")
// 	if err != nil {
// 		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
// 	}

// 	expected := mocks.Policies

// 	for i, v := range res {
// 		if !reflect.DeepEqual(v, &expected[i]) {
// 			t.Errorf("Test failure! res: %v, expected: %v\n", res, expected)
// 		}
// 	}
// }

func initializeAmsMockApi() {
	g := gin.Default()

	g.GET("/users/:agentid", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mocks.Users)
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
