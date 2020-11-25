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

	functionHostName       = "nd-matchfunction.default.svc.cluster.local"
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
		fmt.Printf("Director: fail to get response stream from backend.FetchMatches call, desc: %s\n", err.Error())
		return nil, err
	}

	var matches []*pb.Match
	count := 0
	for {
		resp, err := stream.Recv()
		// assign match
		match := resp.GetMatch()
		matches = append(matches, match)
		fmt.Printf("Got match (Id:%s, Func:%s) - %d Tickets", match.GetMatchId(), match.GetMatchFunction(), len(match.GetTickets()))
		assignErr := assign(be, agonesClient, match)
		if assignErr != nil {
			fmt.Printf("Assign game server failed, got %s", assignErr.Error())
			fmt.Println()
		}
		if err == io.EOF {
			fmt.Println("Resp EOF")
			break
		} else {
			fmt.Printf("Get resp : %d", count)
			fmt.Println()

			count++
		}

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
