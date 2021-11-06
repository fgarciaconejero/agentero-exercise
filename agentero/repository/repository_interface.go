package repository

import "github.com/agentero-exercise/agentero/resources/protos"

type IRepository interface {
	// GetPolicyHoldersFromAms(string) ([]*protos.PolicyHolder, error)
	GetById(string) ([]*protos.PolicyHolder, error)
	GetByMobileNumber(string) (*protos.PolicyHolder, error)
	Upsert([]*protos.PolicyHolder) error
}
