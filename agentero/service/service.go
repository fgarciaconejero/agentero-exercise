package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/agentero-exercise/agentero/domain/models"
	"github.com/agentero-exercise/agentero/repository"
	"github.com/agentero-exercise/agentero/resources/protos"
)

type Service struct {
	repository repository.IRepository
}

func NewService(r repository.IRepository) *Service {
	return &Service{repository: r}
}

func (*Service) GetInsuranceAgentsFromAms() (agents []*models.Agent, err error) {
	resp, err := http.Get("http://localhost:8081/agents/")
	if err != nil {
		log.Fatalln("Error trying to GET '/agents' : ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error trying to read '/agents' response's body: ", err)
	}

	err = json.Unmarshal(body, &agents)
	if err != nil {
		log.Fatalln("Error trying to Unmarshal response from '/agents': ", err)
	}

	return
}

func (s *Service) UpsertInsuranceAgentsIntoSQLite(agents []*models.Agent) (err error) {
	for _, v := range agents {
		err = s.repository.UpsertInsuranceAgent(v)
	}
	return
}

func (*Service) GetPolicyHoldersFromAms(agentId string) (policyHolders []*protos.PolicyHolder, err error) {
	resp, err := http.Get("http://localhost:8081/users/" + agentId)
	if err != nil {
		fmt.Println("There was an unexpected error:", err)
		return nil, err
	} else if resp.StatusCode != 200 {
		fmt.Println("Status code was:", resp.StatusCode)
		return nil, errors.New("HTTP " + fmt.Sprint(resp.StatusCode) + ": Bad Request")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("There was en error when trying to read the response body:", err)
		return
	}

	err = json.Unmarshal(body, &policyHolders)
	if err != nil {
		fmt.Println("There was en error when trying to unmarshal the response body:", err)
		return
	}

	return
}

func (*Service) GetInsurancePoliciesFromAms(agentId string) (insurancePolicies []*protos.InsurancePolicy, err error) {
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

	err = json.Unmarshal(body, &insurancePolicies)
	if err != nil {
		log.Fatalln(err)
	}

	return
}

func (s *Service) GetAllInsuranceAgentsIds() ([]string, error) {
	return s.repository.GetAllInsuranceAgentsIds()
}

func (s *Service) GetContactAndPoliciesByIdFromSQLite(id string) ([]*protos.PolicyHolder, error) {
	return s.repository.GetById(id)
}

func (s *Service) GetContactAndPoliciesByMobileNumberFromSQLite(mobileNumber string) (*protos.PolicyHolder, error) {
	return s.repository.GetByMobileNumber(mobileNumber)
}

func (s *Service) UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite(phs []*protos.PolicyHolder, agentId string) (err error) {
	for _, ph := range phs {
		err = s.repository.UpsertPolicyHolder(ph)
		if err != nil {
			return
		}

		for _, ip := range ph.InsurancePolicy {
			err = s.repository.UpsertInsurancePolicy(ip, agentId)
			if err != nil {
				return
			}
		}
	}
	return
}
