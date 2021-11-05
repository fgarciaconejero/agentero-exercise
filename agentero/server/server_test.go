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
var getByMobileNumberTestingParameters = []struct {
	name         string
	mobileNumber string
	expected     *protos.GetContactsAndPoliciesByMobileNumberResponse
	err          error
}{
	{
		// TODO: Something is failing here, when using 42 this should fail as the one below
		"successful",
		"42",
		&protos.GetContactsAndPoliciesByMobileNumberResponse{
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
		},
		nil,
	},
	{
		"policy holder not found",
		"42",
		nil,
		errors.New("policy holder not found"),
	},
}

func GetContactsAndPoliciesByMobileNumberTest(t *testing.T) {
	for _, tt := range getByMobileNumberTestingParameters {
		s := NewServer(&mockService{})
		req := protos.GetContactsAndPoliciesByMobileNumberRequest{
			MobileNumber: tt.mobileNumber,
		}

		res, err := s.GetContactsAndPoliciesByMobileNumber(context.Background(), &req)
		if err != tt.err {
			t.Errorf("Test failure! res: %v, err: %v\n", res, err)
		}

		if !reflect.DeepEqual(res, tt.expected) {
			t.Errorf("Test failure! res: %v,\n expected: %v\n", res, tt.expected)
		}
	}
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
