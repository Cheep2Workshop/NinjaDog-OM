package main

import (
	"fmt"
	"log"

	"agones.dev/agones/pkg/client/clientset/versioned"
	"google.golang.org/grpc"
	"k8s.io/client-go/rest"
	"open-match.dev/open-match/pkg/pb"
)

func createOMBackendClient() (pb.BackendServiceClient, func() error) {
	conn, err := grpc.Dial(omBackendEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Open Match Backend, got %s", err.Error())
	}
	fmt.Println("Succeed to connect to Open Match Backend")
	return pb.NewBackendServiceClient(conn), conn.Close
}

func createAgonesClient() *versioned.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to connect agones service, go %s", err.Error())
	}

	agonesClient, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to new for config, go %s", err.Error())
	}

	fmt.Println("Succeed to connect to Agones")

	testClientGo()

	return agonesClient
}
