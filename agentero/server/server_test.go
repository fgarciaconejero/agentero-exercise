package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/agentero-exercise/agentero/domain/models"
	"github.com/agentero-exercise/agentero/resources/protos"
)

var getByIdTestingParameters = []struct {
	name     string
	id       string
	expected *protos.GetContactAndPoliciesByIdResponse
	service  mockService
	err      error
}{
	{
		"successful",
		"1",
		&protos.GetContactAndPoliciesByIdResponse{
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
				{
					Name:         "Mary",
					MobileNumber: "000000002",
					InsurancePolicy: []*protos.InsurancePolicy{
						{
							MobileNumber: "000000002",
							Premium:      500,
							Type:         "homeowner",
						},
					},
				},
			},
		},
		mockService{isError: false},
		nil,
	},
	{
		"policy holder not found",
		"2",
		nil,
		mockService{isError: true},
		errors.New("policy holder not found"),
	},
}

func TestGetContactsAndPoliciesById(t *testing.T) {
	for _, tt := range getByIdTestingParameters {
		s := NewServer(&tt.service)
		req := protos.GetContactAndPoliciesByIdRequest{
			InsuranceAgentId: tt.id,
		}
		res, err := s.GetContactAndPoliciesById(context.Background(), &req)
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

var getByMobileNumberTestingParameters = []struct {
	name         string
	mobileNumber string
	expected     *protos.GetContactsAndPoliciesByMobileNumberResponse
	service      mockService
	err          error
}{
	{
		"successful",
		"000000001",
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
		mockService{isError: false},
		nil,
	},
	{
		"policy holder not found",
		"000000002",
		nil,
		mockService{isError: true},
		errors.New("policy holder not found"),
	},
}

func TestGetContactsAndPoliciesByMobileNumber(t *testing.T) {
	for _, tt := range getByMobileNumberTestingParameters {
		s := NewServer(&tt.service)
		req := protos.GetContactsAndPoliciesByMobileNumberRequest{
			MobileNumber: tt.mobileNumber,
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

type mockService struct {
	isError bool
}

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

func (s *mockService) GetContactAndPoliciesByIdFromSQLite(id string) ([]*protos.PolicyHolder, error) {
	if s.isError {
		return nil, errors.New("policy holder not found")
	}
	return []*protos.PolicyHolder{
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
		{
			Name:         "Mary",
			MobileNumber: "000000002",
			InsurancePolicy: []*protos.InsurancePolicy{
				{
					MobileNumber: "000000002",
					Premium:      500,
					Type:         "homeowner",
				},
			},
		},
	}, nil
}

func (s *mockService) GetContactAndPoliciesByMobileNumberFromSQLite(mobileNumber string) (*protos.PolicyHolder, error) {
	if s.isError {
		return nil, errors.New("policy holder not found")
	}
	return &protos.PolicyHolder{
		Name:         "John",
		MobileNumber: "000000001",
		InsurancePolicy: []*protos.InsurancePolicy{
			{
				MobileNumber: "000000001",
				Premium:      500,
				Type:         "homeowner",
			},
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

func (*mockService) UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite(phs []*protos.PolicyHolder, agentId string) error {
	return nil
}

func (*mockService) GetInsuranceAgentsFromAms() (agents []*models.Agent, err error) {
	return
}

func (*mockService) UpsertInsuranceAgentsIntoSQLite(agents []*models.Agent) error {
	return nil
}
