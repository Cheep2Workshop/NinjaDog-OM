package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"agones.dev/agones/pkg/client/clientset/versioned"
	"open-match.dev/open-match/pkg/pb"
)

const (
	omBackendEndpoint = "om-backend.open-match.svc.cluster.local:50505"

	functionHostName       = "nd-matchfunction.nd-om-components.svc.cluster.local"
	functionPort     int32 = 50502
)

func main() {
	for {
		if err := run(); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func run() error {
	bc, closer := createOMBackendClient()
	defer closer()

	agonesClient := createAgonesClient()

	//p := generateMatchProfile()
	matches, err := fetch(bc, agonesClient)
	if err != nil {
		return err
	}
	log.Printf("Generated %v matches", len(matches))

	time.Sleep(time.Second * 5)
	return nil

}

func fetch(be pb.BackendServiceClient, agonesClient *versioned.Clientset) ([]*pb.Match, error) {
	req := createOMFetchMatchesRequest()

	stream, err := be.FetchMatches(context.Background(), req)
	if err != nil {
		log.Println()
		return nil, err
	}

	var matches []*pb.Match
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		match := resp.GetMatch()
		assign(be, agonesClient, match)
		matches = append(matches, match)
	}
	return matches, nil
}

func createOMFetchMatchesRequest() *pb.FetchMatchesRequest {

	return &pb.FetchMatchesRequest{
		// om-function:50502 -> the internal hostname & port number of the MMF service in our Kubernetes cluster
		Config: &pb.FunctionConfig{
			Host: functionHostName,
			Port: functionPort,
			Type: pb.FunctionConfig_GRPC,
		},
		Profile: generateMatchProfile(),
	}

}

func generateMatchProfile() *pb.MatchProfile {
	mode := "mode.demo"
	mp := &pb.MatchProfile{
		Name: "mode_based_profile",
		Pools: []*pb.Pool{
			{
				Name: "pool_mode_" + mode,
				TagPresentFilters: []*pb.TagPresentFilter{
					{
						Tag: mode,
					},
				},
			},
		},
	}
	return mp
}
