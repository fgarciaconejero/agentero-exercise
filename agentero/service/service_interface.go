package service

import (
	"github.com/agentero-exercise/agentero/domain/models"
	"github.com/agentero-exercise/agentero/resources/protos"
)

type IService interface {
	GetInsuranceAgentsFromAms() ([]*models.Agent, error)
	GetPolicyHoldersFromAms(string) ([]*protos.PolicyHolder, error)
	GetInsurancePoliciesFromAms(string) ([]*protos.InsurancePolicy, error)
	GetContactAndPoliciesByIdFromSQLite(string) ([]*protos.PolicyHolder, error)
	// TODO: GetInsurancePoliciesFromSQLite
	GetAllInsuranceAgentsIds() ([]string, error)
	UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite([]*protos.PolicyHolder, string) error
	UpsertInsuranceAgentsIntoSQLite([]*models.Agent) error
}
