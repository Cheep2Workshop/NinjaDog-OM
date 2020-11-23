package main

import "fmt"

// import (
// 	"log"

// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/rest"
// )

// func testClientGo() {
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		log.Fatalf("Cannot get config, got %s", err.Error())
// 	}

// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		log.Fatalf("Cannot refresh config, got %s", err.Error())
// 	}

// 	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
// 	if err != nil {
// 		log.Fatalf("Cannot list pods, got %s", err.Error())
// 	}

// }

func testClientGo() {
	fmt.Println("Test client go")
}
