package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/agentero-exercise/agentero/repository"
	"github.com/agentero-exercise/agentero/resources/protos"
	"github.com/agentero-exercise/agentero/service"
	"google.golang.org/grpc"
)

type server struct {
	Service service.IService
}

func NewServer(s service.IService) *server {
	return &server{Service: s}
}

// TODO: Add logs
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n\n", err)
	}

	r, err := repository.NewRepository()
	if err != nil {
		log.Fatalln("There was an error while creating a new repository: ", err)
	}
	srv := NewServer(service.NewService(r))
	s := grpc.NewServer()

	protos.RegisterPolicyHoldersServiceServer(s, srv)
	fmt.Println("Created server successfully!")

	agents, err := srv.Service.GetInsuranceAgentsFromAms()
	if err != nil {
		log.Fatalf("There was an unexpected error on GetInsuranceAgentsFromAms: %v\n", err)
	}
	fmt.Println("Got insurance agents from AMS")

	err = srv.Service.UpsertInsuranceAgentsIntoSQLite(agents)
	if err != nil {
		log.Fatalf("There was an unexpected error on UpsertInsuranceAgentsIntoSQLite: %v\n", err)
	}
	fmt.Println("Upsert of Insurance Agents successful!")

	agentIds, err := srv.Service.GetAllInsuranceAgentsIds()
	if err != nil {
		log.Fatalf("There was an unexpected error on GetAllInsuranceAgentsIds: %v\n", err)
	}
	fmt.Println("Got insurance agents ids from SQLite")

	for _, id := range agentIds {
		phs, err := srv.GetPolicyHoldersAndInsurancePoliciesFromAms(id)
		if err != nil {
			log.Fatalf("There was an unexpected error on GetPolicyHoldersFromAms: %v\n", err)
		}

		err = srv.Service.UpsertPolicyHoldersAndInsurancePoliciesIntoSQLite(phs, id)
		if err != nil {
			log.Fatalf("There was an unexpected error while trying to Upsert to SQLite: %v\n", err)
		}
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
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

func (s *server) GetContactsAndPoliciesByMobileNumber(ctx context.Context, req *protos.GetContactsAndPoliciesByMobileNumberRequest) (*protos.GetContactsAndPoliciesByMobileNumberResponse, error) {
	res, err := s.Service.GetContactAndPoliciesByMobileNumberFromSQLite(req.MobileNumber)
	if err != nil {
		return nil, err
	}

	return &protos.GetContactsAndPoliciesByMobileNumberResponse{
		PolicyHolder: res,
	}, nil
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
