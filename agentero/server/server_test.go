package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/stretchr/testify/assert"
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
				MobileNumber: "43",
				InsurancePolicy: []*protos.InsurancePolicy{
					{
						MobileNumber: "43",
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

// TODO: turn this into parameterized tests
func TestGetContactsAndPoliciesByMobileNumber_Success(t *testing.T) {
	s := NewServer(&mockService{})
	req := protos.GetContactsAndPoliciesByMobileNumberRequest{
		MobileNumber: "43",
	}

	res, err := s.GetContactsAndPoliciesByMobileNumber(context.Background(), &req)
	if err != nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}

	expected := &protos.GetContactsAndPoliciesByMobileNumberResponse{
		PolicyHolder: &protos.PolicyHolder{
			Name:         "John",
			MobileNumber: "43",
			InsurancePolicy: []*protos.InsurancePolicy{
				{
					MobileNumber: "43",
					Premium:      500,
					Type:         "homeowner",
				},
			},
		},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Test failure! res: %v,\n expected: %v\n", res, expected)
	}
}

func TestGetContactsAndPoliciesByMobileNumber_PolicyHolderNotFound(t *testing.T) {
	s := NewServer(&mockService{})
	req := protos.GetContactsAndPoliciesByMobileNumberRequest{
		MobileNumber: "42",
	}

	res, err := s.GetContactsAndPoliciesByMobileNumber(context.Background(), &req)
	if err == nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}

	assert.Equal(t,
		errors.New("policy holder not found"), err)

}

type mockService struct{}

func (*mockService) GetPolicyHoldersFromAms(agentId string) ([]*protos.PolicyHolder, error) {
	return []*protos.PolicyHolder{
		{
			Name:         "John",
			MobileNumber: "43",
		},
	}, nil
}

func (*mockService) GetInsurancePoliciesFromAms(agentId string) ([]*protos.InsurancePolicy, error) {
	return []*protos.InsurancePolicy{
		{
			MobileNumber: "43",
			Premium:      500,
			Type:         "homeowner",
		},
	}, nil
}
