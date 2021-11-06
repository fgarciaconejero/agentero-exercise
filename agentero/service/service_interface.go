package service

import "github.com/agentero-exercise/agentero/resources/protos"

type IService interface {
	GetPolicyHoldersFromAms(string) ([]*protos.PolicyHolder, error)
	GetInsurancePoliciesFromAms(string) ([]*protos.InsurancePolicy, error)
	GetAllInsuranceAgentsIds() ([]string, error)
	UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite([]*protos.PolicyHolder) error
}
