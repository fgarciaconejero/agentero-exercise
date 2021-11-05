package main

import (
	"context"
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

func TestGetContactsAndPoliciesByMobileNumber(t *testing.T) {
	s := NewServer(&mockService{})
	req := protos.GetContactsAndPoliciesByMobileNumberRequest{
		MobileNumber: "some-mobile-number",
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
