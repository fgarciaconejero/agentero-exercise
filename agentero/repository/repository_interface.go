package repository

import "github.com/agentero-exercise/agentero/resources/protos"

type IRepository interface {
	// GetPolicyHoldersFromAms(string) ([]*protos.PolicyHolder, error)
	Get(string) (*protos.PolicyHolder, error)
	Upsert(*protos.PolicyHolder) error
}
