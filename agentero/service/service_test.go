package service_test

import (
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/agentero-exercise/agentero/resources/mocks"
	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/agentero-exercise/agentero/service"
	"github.com/gin-gonic/gin"
)

func TestGetPolicyHoldersFromAms(t *testing.T) {
	initializeAmsMockApi()

	s := service.NewService(&mockRepository{})

	res, err := s.GetPolicyHoldersFromAms("some-agent-id")
	if err != nil {
		t.Errorf("Error failure! res: %v,\n err: %v\n", res, err)
	}

	expected := mocks.Users

	for i, v := range res {
		if !reflect.DeepEqual(v, &expected[i]) {
			t.Errorf("Mismatch failure! res: %v,\n\n\n expected: %v\n", v, &expected[i])
		}
	}
}

func GetInsurancePoliciesFromAms(t *testing.T) {
	initializeAmsMockApi()

	s := service.NewService(&mockRepository{})

	res, err := s.GetInsurancePoliciesFromAms("some-agent-id")
	if err != nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}

	expected := mocks.Policies

	for i, v := range res {
		if !reflect.DeepEqual(v, &expected[i]) {
			t.Errorf("Test failure! res: %v, expected: %v\n", res, expected)
		}
	}
}

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

func (*mockRepository) UpsertInsurancePolicy(phs *protos.InsurancePolicy) error {
	return nil
}

func (r *mockRepository) GetAllInsuranceAgentsIds() ([]string, error) {
	return []string{}, nil
}
