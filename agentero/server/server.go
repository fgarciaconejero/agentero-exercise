package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/agentero-exercise/agentero/resources/protos"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n\n", err)
	}

	s := grpc.NewServer()
	protos.RegisterPolicyHoldersServiceServer(s, &server{})
	fmt.Printf("Created server: %v\n", s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

type server struct{}

func (*server) GetContactAndPoliciesById(ctx context.Context, req *protos.GetContactAndPoliciesByIdRequest) (*protos.GetContactAndPoliciesByIdResponse, error) {
	resp, err := http.Get("http://localhost:8081/users/1")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	x := []*protos.PolicyHolder{}
	err = json.Unmarshal(body, &x)
	if err != nil {
		log.Fatalln(err)
	}
	// Adding this log for test purposes only
	fmt.Printf("%v", x)

	return &protos.GetContactAndPoliciesByIdResponse{
		PolicyHolders: x,
	}, nil
}

func (*server) GetContactsAndPoliciesByMobileNumber(ctx context.Context, req *protos.GetContactsAndPoliciesByMobileNumberRequest) (*protos.GetContactsAndPoliciesByMobileNumberResponse, error) {
	return nil, nil
}
