package main

import (
	"context"
	"testing"

	"github.com/agentero-exercise/agentero/resources/protos"
)

func TestGetContactAndPoliciesById(t *testing.T) {

	s := &server{}
	req := protos.GetContactAndPoliciesByIdRequest{
		InsuranceAgentId: "some-id",
	}

	res, err := s.GetContactAndPoliciesById(context.Background(), &req)
	if err != nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}
}

func TestGetContactsAndPoliciesByMobileNumber(t *testing.T) {
	s := &server{}
	req := protos.GetContactsAndPoliciesByMobileNumberRequest{
		MobileNumber: "some-mobile-number",
	}
	res, err := s.GetContactsAndPoliciesByMobileNumber(context.Background(), &req)
	if res != nil || err != nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}
}
