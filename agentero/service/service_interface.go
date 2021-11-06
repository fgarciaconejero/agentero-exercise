package service

import "github.com/agentero-exercise/agentero/resources/protos"

type IService interface {
	GetPolicyHoldersFromAms(string) ([]*protos.PolicyHolder, error)
	GetInsurancePoliciesFromAms(string) ([]*protos.InsurancePolicy, error)
	GetPolicyHoldersFromSQLite(string) ([]*protos.PolicyHolder, error)
	// TODO: GetInsurancePoliciesFromSQLite
	GetAllInsuranceAgentsIds() ([]string, error)
	UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite([]*protos.PolicyHolder) error
}
