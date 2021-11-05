package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/agentero-exercise/agentero/resources/protos"
)

type Service struct{}

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
