package main

import (
	"context"
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

	agentIds, err := srv.Service.GetAllInsuranceAgentsIds()
	if err != nil {
		log.Fatalf("There was an unexpected error on GetAllInsuranceAgentsIds: %v\n", err)
	}

	for _, id := range agentIds {
		// TODO: Retrieve all agentIDs from db to get all from AMS
		phs, err := srv.GetPolicyHoldersAndInsurancePoliciesFromAms(id)
		if err != nil {
			log.Fatalf("There was an unexpected error on GetPolicyHoldersFromAms: %v\n", err)
		}

		err = srv.Service.UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite(phs)
		if err != nil {
			log.Fatalf("There was an unexpected error while trying to Upsert to SQLite: %v\n", err)
		}

	}

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

func (s *server) GetPolicyHoldersAndInsurancePoliciesFromAms(id string) ([]*protos.PolicyHolder, error) {
	phs, err := s.Service.GetPolicyHoldersFromAms(id)
	if err != nil {
		return nil, err
	}

	ips, err := s.Service.GetInsurancePoliciesFromAms(id)
	if err != nil {
		return nil, err
	}

	filterMobileNumbers(phs, ips)
	mapPoliciesToHolders(ips, phs)

	return phs, nil
}

func (s *server) GetContactAndPoliciesById(ctx context.Context, req *protos.GetContactAndPoliciesByIdRequest) (*protos.GetContactAndPoliciesByIdResponse, error) {
	res, err := s.Service.GetContactAndPoliciesByIdFromSQLite(req.InsuranceAgentId)
	if err != nil {
		return nil, err
	}

	return &protos.GetContactAndPoliciesByIdResponse{
		PolicyHolders: res,
	}, nil
}

// TODO: This should retrieve from db, not ams
func (s *server) GetContactsAndPoliciesByMobileNumber(ctx context.Context, req *protos.GetContactsAndPoliciesByMobileNumberRequest) (*protos.GetContactsAndPoliciesByMobileNumberResponse, error) {
	return &protos.GetContactsAndPoliciesByMobileNumberResponse{}, nil
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
