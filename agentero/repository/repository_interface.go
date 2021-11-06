package repository

import "github.com/agentero-exercise/agentero/resources/protos"

type IRepository interface {
	GetById(string) ([]*protos.PolicyHolder, error)
	GetByMobileNumber(string) (*protos.PolicyHolder, error)
	GetAllInsuranceAgentsIds() (string, error)
	UpsertPolicyHolder(*protos.PolicyHolder) error
	UpsertInsurancePolicy(*protos.InsurancePolicy) error
}
