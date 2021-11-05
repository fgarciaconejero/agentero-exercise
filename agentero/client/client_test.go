package main

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/agentero-go/policy_holder/policy_holder_pb"
	"google.golang.org/grpc"
)

func TestGetById(t *testing.T) {
	newMockServer()
}

func newMockServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n\n", err)
	}

	s := grpc.NewServer()
	policy_holder_pb.RegisterPolicyHoldersServiceServer(s, &mockServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

type mockServer struct {
}

func (*mockServer) GetContactAndPoliciesById(ctx context.Context, req *policy_holder_pb.GetContactAndPoliciesByIdRequest) (*policy_holder_pb.GetContactAndPoliciesByIdResponse, error) {
	return nil, nil
}

func (*mockServer) GetContactsAndPoliciesByMobileNumber(ctx context.Context, req *policy_holder_pb.GetContactsAndPoliciesByMobileNumberRequest) (*policy_holder_pb.GetContactsAndPoliciesByMobileNumberResponse, error) {
	return nil, nil
}
