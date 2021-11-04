package main

import (
	"context"
	"testing"

	"github.com/agentero-go/policy_holder/policy_holder_pb"
)

func TestGetContactAndPoliciesById(t *testing.T) {
	s := &server{}
	req := policy_holder_pb.GetContactAndPoliciesByIdRequest{
		InsuranceAgentId: "some-id",
	}
	res, err := s.GetContactAndPoliciesById(context.Background(), &req)
	if res != nil || err != nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}
}

func TestGetContactsAndPoliciesByMobileNumber(t *testing.T) {
	s := &server{}
	req := policy_holder_pb.GetContactsAndPoliciesByMobileNumberRequest{
		MobileNumber: "some-mobile-number",
	}
	res, err := s.GetContactsAndPoliciesByMobileNumber(context.Background(), &req)
	if res != nil || err != nil {
		t.Errorf("Test failure! res: %v, err: %v\n", res, err)
	}
}
