package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"

	"github.com/agentero-exercise/agentero/repository"
	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/agentero-exercise/agentero/service"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n\n", err)
	}

	r, err := repository.NewRepository()
	if err != nil {
		log.Fatalln("There was an error while creating a new server: ", err)
	}
	srv := NewServer(service.NewService(r))
	s := grpc.NewServer()

	protos.RegisterPolicyHoldersServiceServer(s, srv)
	fmt.Println("Created server successfuly!")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

type server struct {
	Service service.IService
}

func NewServer(s service.IService) *server {
	return &server{Service: s}
}

func (s *server) GetContactAndPoliciesById(ctx context.Context, req *protos.GetContactAndPoliciesByIdRequest) (*protos.GetContactAndPoliciesByIdResponse, error) {
	phs, ips, err := s.getPolicyHoldersAndInsurancePoliciesFromAms(req.InsuranceAgentId)
	if err != nil {
		// Not logging anything since it's already being logged in the function
		return nil, err
	}

	filterMobileNumbers(phs, ips)
	mapPoliciesToHolders(ips, phs)

	// TODO: Save the data in-memory or a database before returning it
	return &protos.GetContactAndPoliciesByIdResponse{
		PolicyHolders: phs,
	}, nil
}

func (s *server) GetContactsAndPoliciesByMobileNumber(ctx context.Context, req *protos.GetContactsAndPoliciesByMobileNumberRequest) (*protos.GetContactsAndPoliciesByMobileNumberResponse, error) {
	phs, ips, err := s.getPolicyHoldersAndInsurancePoliciesFromAms(req.InsuranceAgentId)
	if err != nil {
		// Not logging anything since it's already being logged in the function
		return nil, err
	}

	filterMobileNumbers(phs, ips)
	mapPoliciesToHolders(ips, phs)

	ph, err := filterPolicyHolderByMobileNumber(phs, req.MobileNumber)
	if err != nil {
		return nil, err
	}

	return &protos.GetContactsAndPoliciesByMobileNumberResponse{
		PolicyHolder: ph,
	}, nil
}

func (s *server) getPolicyHoldersAndInsurancePoliciesFromAms(id string) ([]*protos.PolicyHolder, []*protos.InsurancePolicy, error) {
	phs, err := s.Service.GetPolicyHoldersFromAms(id)
	if err != nil {
		log.Fatalf("There was an unexpected error on GetPolicyHoldersFromAms: %v\n", err)
		return nil, nil, err
	}

	ips, err := s.Service.GetInsurancePoliciesFromAms(id)
	if err != nil {
		log.Fatalf("There was an unexpected error on GetInsurancePoliciesFromAms: %v\n", err)
		return nil, nil, err
	}

	return phs, ips, nil
}

// Returns the first policy holder whose mobile number matches the desired one, otherwise it returns an error
func filterPolicyHolderByMobileNumber(phs []*protos.PolicyHolder, mobileNumber string) (*protos.PolicyHolder, error) {
	for _, v := range phs {
		if v.MobileNumber == mobileNumber {
			return v, nil
		}
	}
	return nil, errors.New("policy holder not found")
}

// Removes every character that is not a number from Mobile Numbers of both Insurance Policies and Policy Holders
func filterMobileNumbers(phs []*protos.PolicyHolder, ips []*protos.InsurancePolicy) *regexp.Regexp {
	// This regexp filters everything but numbers out
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	formatMobileNumbersFromInsurancePolicies(ips, reg)
	formatMobileNumbersFromPolicyHolders(phs, reg)

	return reg
}

// Inserts insurance policies into their rightful policy holders
func mapPoliciesToHolders(ips []*protos.InsurancePolicy, phs []*protos.PolicyHolder) {
	for _, iPolicy := range ips {
		for _, pHolder := range phs {
			if iPolicy.MobileNumber == pHolder.MobileNumber {
				pHolder.InsurancePolicy = append(pHolder.InsurancePolicy, iPolicy)
			}
		}
	}
}

// Removes every character that is not a number from Mobile Numbers of Insurance Policies
func formatMobileNumbersFromInsurancePolicies(ips []*protos.InsurancePolicy, reg *regexp.Regexp) {
	for i, v := range ips {
		ips[i].MobileNumber = reg.ReplaceAllString(v.MobileNumber, "")
	}
}

// Removes every character that is not a number from Mobile Numbers of Policy Holders
func formatMobileNumbersFromPolicyHolders(phs []*protos.PolicyHolder, reg *regexp.Regexp) {
	for i, v := range phs {
		phs[i].MobileNumber = reg.ReplaceAllString(v.MobileNumber, "")
	}
}
