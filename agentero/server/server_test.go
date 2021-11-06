package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/agentero-exercise/agentero/resources/protos"
)

func TestGetContactAndPoliciesById(t *testing.T) {

	s := NewServer(&mockService{})
	req := protos.GetContactAndPoliciesByIdRequest{
		InsuranceAgentId: "some-id",
	}

	res, err := s.GetContactAndPoliciesById(context.Background(), &req)
	if err != nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}

	expected := &protos.GetContactAndPoliciesByIdResponse{
		PolicyHolders: []*protos.PolicyHolder{
			{
				Name:         "John",
				MobileNumber: "000000001",
				InsurancePolicy: []*protos.InsurancePolicy{
					{
						MobileNumber: "000000001",
						Premium:      500,
						Type:         "homeowner",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Test failure! res: %v,\n expected: %v\n", res, expected)
	}
}

var getByMobileNumberTestingParameters = []struct {
	name             string
	mobileNumber     string
	insuranceAgentId string
	expected         *protos.GetContactsAndPoliciesByMobileNumberResponse
	err              error
}{
	{
		"successful",
		"000000001",
		"some-agent-id",
		&protos.GetContactsAndPoliciesByMobileNumberResponse{
			PolicyHolder: &protos.PolicyHolder{
				Name:         "John",
				MobileNumber: "000000001",
				InsurancePolicy: []*protos.InsurancePolicy{
					{
						MobileNumber: "000000001",
						Premium:      500,
						Type:         "homeowner",
					},
				},
			},
		},
		nil,
	},
	{
		"policy holder not found",
		"000000002",
		"some-agent-id",
		nil,
		errors.New("policy holder not found"),
	},
}

func TestGetContactsAndPoliciesByMobileNumber(t *testing.T) {
	for _, tt := range getByMobileNumberTestingParameters {
		s := NewServer(&mockService{})
		req := protos.GetContactsAndPoliciesByMobileNumberRequest{
			MobileNumber:     tt.mobileNumber,
			InsuranceAgentId: tt.insuranceAgentId,
		}

		res, err := s.GetContactsAndPoliciesByMobileNumber(context.Background(), &req)
		// Lint warns not to use DeepEqual on error, but every other way doesn't work or panics because
		// in the case of the error being nil there is a nil pointer exception
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Test '%v' failed! err: %v, expected: %v\n", tt.name, err, tt.err)
		}

		if !reflect.DeepEqual(res, tt.expected) {
			t.Errorf("Test '%v' failed! res: %v,\n expected: %v\n", tt.name, res, tt.expected)
		}
	}
}

type mockService struct{}

func (s *mockService) GetAllInsuranceAgentsIds() ([]string, error) {
	return nil, nil
}

func (*mockService) GetPolicyHoldersFromAms(agentId string) ([]*protos.PolicyHolder, error) {
	return []*protos.PolicyHolder{
		{
			Name:         "John",
			MobileNumber: "000000001",
		},
	}, nil
}

func (*mockService) GetInsurancePoliciesFromAms(agentId string) ([]*protos.InsurancePolicy, error) {
	return []*protos.InsurancePolicy{
		{
			MobileNumber: "000000001",
			Premium:      500,
			Type:         "homeowner",
		},
	}, nil
}

func (*mockService) UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite(phs []*protos.PolicyHolder) error {
	return nil
}
