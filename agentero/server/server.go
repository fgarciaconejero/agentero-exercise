package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/agentero-go/policy_holder/policy_holder_pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n\n", err)
	}

	s := grpc.NewServer()
	policy_holder_pb.RegisterPolicyHoldersServiceServer(s, &server{})
	fmt.Printf("Created server: %v\n", s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

type server struct{}

func (*server) GetContactAndPoliciesById(ctx context.Context, req *policy_holder_pb.GetContactAndPoliciesByIdRequest) (*policy_holder_pb.GetContactAndPoliciesByIdResponse, error) {
	return nil, nil
}

func (*server) GetContactsAndPoliciesByMobileNumber(ctx context.Context, req *policy_holder_pb.GetContactsAndPoliciesByMobileNumberRequest) (*policy_holder_pb.GetContactsAndPoliciesByMobileNumberResponse, error) {
	return nil, nil
}
