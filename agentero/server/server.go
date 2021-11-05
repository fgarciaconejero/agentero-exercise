package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"regexp"

	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/agentero-exercise/agentero/service"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n\n", err)
	}

	s := grpc.NewServer()
	protos.RegisterPolicyHoldersServiceServer(s, NewServer(&service.Service{}))
	fmt.Println("Created server successfuly!")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

type server struct {
	Service service.IService
}

func NewServer(service service.IService) *server {
	return &server{Service: service}
}

func (s *server) GetContactAndPoliciesById(ctx context.Context, req *protos.GetContactAndPoliciesByIdRequest) (*protos.GetContactAndPoliciesByIdResponse, error) {
	phs, err := s.Service.GetPolicyHoldersFromAms(req.InsuranceAgentId)
	if err != nil {
		log.Fatalf("There was an unexpected error on GetPolicyHoldersFromAms: %v\n", err)
	}

	ips, err := s.Service.GetInsurancePoliciesFromAms(req.InsuranceAgentId)
	if err != nil {
		log.Fatalf("There was an unexpected error on GetInsurancePoliciesFromAms: %v\n", err)
	}

	// Remove every character that is not a number from Mobile Numbers
	filterMobileNumberRegexp := filterMobileNumberRegexp()
	formatMobileNumbersFromInsurancePolicies(ips, filterMobileNumberRegexp)
	formatMobileNumbersFromPolicyHolders(phs, filterMobileNumberRegexp)

	mapPoliciesToHolders(ips, phs)

	return &protos.GetContactAndPoliciesByIdResponse{
		PolicyHolders: phs,
	}, nil
}

func (*server) GetContactsAndPoliciesByMobileNumber(ctx context.Context, req *protos.GetContactsAndPoliciesByMobileNumberRequest) (*protos.GetContactsAndPoliciesByMobileNumberResponse, error) {
	return nil, nil
}

func mapPoliciesToHolders(ips []*protos.InsurancePolicy, phs []*protos.PolicyHolder) {
	for _, iPolicy := range ips {
		for _, pHolder := range phs {
			if iPolicy.MobileNumber == pHolder.MobileNumber {
				pHolder.InsurancePolicy = append(pHolder.InsurancePolicy, iPolicy)
			}
		}
	}
}

func formatMobileNumbersFromInsurancePolicies(ips []*protos.InsurancePolicy, reg *regexp.Regexp) {
	for i, v := range ips {
		ips[i].MobileNumber = reg.ReplaceAllString(v.MobileNumber, "")
	}
}

func formatMobileNumbersFromPolicyHolders(phs []*protos.PolicyHolder, reg *regexp.Regexp) {
	for i, v := range phs {
		phs[i].MobileNumber = reg.ReplaceAllString(v.MobileNumber, "")
	}
}

func filterMobileNumberRegexp() *regexp.Regexp {
	// This regexp filters everything but numbers out
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	return reg
}
