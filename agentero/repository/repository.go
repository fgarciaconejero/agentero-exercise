package repository

import "github.com/agentero-exercise/agentero/resources/protos"

type Repository struct {
}

func (*Repository) GetById(agentId string) ([]*protos.PolicyHolder, error) {
	return nil, nil
}

func (*Repository) GetByMobileNumber(agentId string) ([]protos.PolicyHolder, error) {
	return nil, nil
}

func (*Repository) Upsert(phs []*protos.PolicyHolder) error {
	return nil
}
