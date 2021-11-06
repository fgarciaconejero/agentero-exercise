package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/agentero-exercise/agentero/repository"
	"github.com/agentero-exercise/agentero/resources/protos"
)

type Service struct {
	repository repository.IRepository
}

func NewService(r repository.IRepository) *Service {
	return &Service{repository: r}
}

func (*Service) GetPolicyHoldersFromAms(agentId string) ([]*protos.PolicyHolder, error) {
	resp, err := http.Get("http://localhost:8081/users/" + agentId)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	policyHolders := []*protos.PolicyHolder{}
	err = json.Unmarshal(body, &policyHolders)
	if err != nil {
		log.Fatalln(err)
	}

	return policyHolders, nil
}

func (*Service) GetInsurancePoliciesFromAms(agentId string) ([]*protos.InsurancePolicy, error) {
	resp, err := http.Get("http://localhost:8081/policies/" + agentId)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	insurancePolicies := []*protos.InsurancePolicy{}
	err = json.Unmarshal(body, &insurancePolicies)
	if err != nil {
		log.Fatalln(err)
	}

	return insurancePolicies, nil
}

func (s *Service) GetAllInsuranceAgentsIds() ([]string, error) {
	return s.GetAllInsuranceAgentsIds()
}

func (s *Service) UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite(phs []*protos.PolicyHolder) error {
	for _, ph := range phs {
		err := s.repository.UpsertPolicyHolder(ph)
		if err != nil {
			return err
		}

		for _, ip := range ph.InsurancePolicy {

			err = s.repository.UpsertInsurancePolicy(ip)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
