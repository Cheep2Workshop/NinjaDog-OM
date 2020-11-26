package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"agones.dev/agones/pkg/client/clientset/versioned"
	"google.golang.org/grpc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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

// https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration
func createAgonesClient() (*versioned.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect agones service, got %s", err.Error())
		// log.Fatalf("Failed to connect agones service, got %s", err.Error())
	}

	agonesClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to new for config, go %s", err.Error())
		// log.Fatalf("Failed to new for config, go %s", err.Error())
	}

	fmt.Println("Succeed to connect to Agones")

	return agonesClient, nil
}

// https://gist.github.com/ks888/0a0e0fbf4694d7955999a6f59aa2766d
func createAgonesClient2() *versioned.Clientset {
	contextName := ""
	if len(os.Args) >= 2 {
		contextName = os.Args[1]
	}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		log.Fatalf("Failed to connect agones service, got %s", err.Error())
	}
	agonesClient, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to new for config, go %s", err.Error())
	}

	fmt.Println("Succeed to connect to Agones")

	return agonesClient
}

// https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go
func createAgonesClient3() *versioned.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	agonesClient, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to new for config, go %s", err.Error())
	}

	fmt.Println("Succeed to connect to Agones")

	return agonesClient
}
