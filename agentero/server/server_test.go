package main

import (
	"context"
	"errors"
	"testing"

	"github.com/agentero-exercise/agentero/domain/models"
	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/stretchr/testify/assert"
)

var getFromAmsTestingParameters = []struct {
	name     string
	id       string
	expected []*protos.PolicyHolder
	service  mockService
	err      error
}{
	{
		"successful",
		"some-agent-id",
		[]*protos.PolicyHolder{
			{
				Name:         "John",
				MobileNumber: "000000001",
				InsurancePolicy: []*protos.InsurancePolicy{
					{
						MobileNumber: "000000001",
						Premium:      500,
						Type:         "homeowner",
						AgentId:      "1",
					},
				},
			},
			{
				Name:         "Mary",
				MobileNumber: "000000002",
				InsurancePolicy: []*protos.InsurancePolicy{
					{
						MobileNumber: "000000002",
						Premium:      20,
						Type:         "homeowner",
						AgentId:      "1",
					},
					{
						MobileNumber: "000000002",
						Premium:      10,
						Type:         "homeowner",
						AgentId:      "1",
					},
				},
			},
		},
		mockService{},
		nil,
	},
	{
		"GetPolicyHoldersFromAms returns an error",
		"some-agent-id",
		nil,
		mockService{isGetPolicyHoldersFromAmsError: true},
		errors.New("GetPolicyHoldersFromAms returned an error"),
	},
	{
		"GetInsurancePoliciesFromAms returns an error",
		"some-agent-id",
		nil,
		mockService{isGetInsurancePoliciesFromAmsError: true},
		errors.New("GetInsurancePoliciesFromAms returned an error"),
	},
}

func TestGetPolicyHoldersAndInsurancePoliciesFromAms(t *testing.T) {
	for _, tt := range getFromAmsTestingParameters {
		s := NewServer(&tt.service)

		res, err := s.GetPolicyHoldersAndInsurancePoliciesFromAms(tt.id)
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
		assert.Equal(t, res, tt.expected)
	}
}

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
							AgentId:      "1",
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
							AgentId:      "1",
						},
					},
				},
			},
		},
		mockService{},
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
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
		assert.Equal(t, res, tt.expected)
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
						AgentId:      "1",
					},
				},
			},
		},
		mockService{},
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
		if tt.err != nil {
			assert.EqualError(t, err, tt.err.Error())
		}
		assert.Equal(t, res, tt.expected)
	}
}

type mockService struct {
	isError                            bool
	isGetPolicyHoldersFromAmsError     bool
	isGetInsurancePoliciesFromAmsError bool
}

func (s *mockService) GetAllInsuranceAgentsIds() ([]string, error) {
	return nil, nil
}

func (s *mockService) GetPolicyHoldersFromAms(agentId string) ([]*protos.PolicyHolder, error) {
	if s.isGetPolicyHoldersFromAmsError {
		return nil, errors.New("GetPolicyHoldersFromAms returned an error")
	}
	return []*protos.PolicyHolder{
		{
			Name:            "John",
			MobileNumber:    "000000001",
			InsurancePolicy: nil,
		},
		{
			Name:            "Mary",
			MobileNumber:    "000000002",
			InsurancePolicy: nil,
		},
	}, nil
}

func (s *mockService) GetInsurancePoliciesFromAms(agentId string) ([]*protos.InsurancePolicy, error) {
	if s.isGetInsurancePoliciesFromAmsError {
		return nil, errors.New("GetInsurancePoliciesFromAms returned an error")
	}
	return []*protos.InsurancePolicy{
		{
			MobileNumber: "000000001",
			Premium:      500,
			Type:         "homeowner",
			AgentId:      "1",
		},
		{
			MobileNumber: "000000002",
			Premium:      20,
			Type:         "homeowner",
			AgentId:      "1",
		},
		{
			MobileNumber: "000000002",
			Premium:      10,
			Type:         "homeowner",
			AgentId:      "1",
		},
	}, nil
}

func (s *mockService) GetInsuranceAgentsFromAms() (agents []*models.Agent, err error) {
	return
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
					AgentId:      "1",
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
					AgentId:      "1",
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
				AgentId:      "1",
			},
		},
	}, nil
}

func (*mockService) UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite(phs []*protos.PolicyHolder, agentId string) error {
	return nil
}

func (*mockService) UpsertInsuranceAgentsIntoSQLite(agents []*models.Agent) error {
	return nil
}
