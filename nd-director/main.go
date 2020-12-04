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

var modes = [2]string{"private", "duel"}

func main() {
	// execute fecth procedure every 1 second
	for range time.Tick(5 * time.Second) {
		err := run()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func run() error {
	bc, closer := createOMBackendClient()
	defer closer()

	agonesClient, err := createAgonesClient()
	if err != nil {
		return err
	}

	for _, mode := range modes {
		matches, err := fetch(bc, agonesClient, mode)
		if len(matches) > 0 {
			log.Printf("Generated %v matches", len(matches))
		}
		if err != nil {
			log.Fatalf("Failed to fetch match, got %s", err.Error())
		}
	}
	return nil
}

func fetch(be pb.BackendServiceClient, agonesClient *versioned.Clientset, mode string) ([]*pb.Match, error) {
	req := createOMFetchMatchesRequest(mode)

	stream, err := be.FetchMatches(context.Background(), req)
	if err != nil {
		fmt.Printf("Director: fail to get response stream from backend.FetchMatches call, desc: %s\n", err.Error())
		fmt.Println()
		return nil, err
	}

	var matches []*pb.Match
	count := 0
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// assign match
		match := resp.GetMatch()
		matches = append(matches, match)
		if match == nil {
			fmt.Println("match is null")
		}
		fmt.Printf("Got match (Id:%s, Func:%s) - %d Tickets", match.GetMatchId(), match.GetMatchFunction(), len(match.GetTickets()))
		fmt.Println()
		assignErr := assign(be, agonesClient, match)
		if assignErr != nil {
			fmt.Printf("Assign game server failed, got %s", assignErr.Error())
			fmt.Println()
		}
		count++
	}
	return matches, nil
}

func createOMFetchMatchesRequest(mode string) *pb.FetchMatchesRequest {

	return &pb.FetchMatchesRequest{
		// om-function:50502 -> the internal hostname & port number of the MMF service in our Kubernetes cluster
		Config: &pb.FunctionConfig{
			Host: functionHostName,
			Port: functionPort,
			Type: pb.FunctionConfig_GRPC,
		},
		Profile: generateMatchProfile(mode),
	}

}

func generateMatchProfile(mode string) *pb.MatchProfile {
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
