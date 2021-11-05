package main

import (
	"context"
	"fmt"
	"log"
	"net"

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
	protos.RegisterPolicyHoldersServiceServer(s, &server{})
	fmt.Println("Created server successfuly!")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

type server struct {
	Service *service.Service
}

func (s *server) GetContactAndPoliciesById(ctx context.Context, req *protos.GetContactAndPoliciesByIdRequest) (*protos.GetContactAndPoliciesByIdResponse, error) {
	s.initializeService()
	phs, err := s.Service.GetPolicyHoldersFromAms(req.InsuranceAgentId)
	if err != nil {
		log.Fatalf("There was an unexpected error on GetPolicyHoldersFromAms: %v\n", err)
	}

	ips, err := s.Service.GetInsurancePoliciesFromAms(req.InsuranceAgentId)
	if err != nil {
		log.Fatalf("There was an unexpected error on GetInsurancePoliciesFromAms: %v\n", err)
	}
	fmt.Println(ips)

	return &protos.GetContactAndPoliciesByIdResponse{
		PolicyHolders: phs,
	}, nil
}

func (*server) GetContactsAndPoliciesByMobileNumber(ctx context.Context, req *protos.GetContactsAndPoliciesByMobileNumberRequest) (*protos.GetContactsAndPoliciesByMobileNumberResponse, error) {
	return nil, nil
}

// This method is needed to be able to mock the service in the tests for the two functions above
func (s *server) initializeService() {
	s.Service = &service.Service{}
}
