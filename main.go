package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Define command-line flags
	var kubeconfig string

	// If the kubeconfig flag is not provided, use the default path
	home := homedir.HomeDir()
	if home == "" {
		fmt.Println("Error: unable to find home directory")
		os.Exit(1)
	}
	defaultKubeconfig := filepath.Join(home, ".kube", "config")

	// Set up the flag for kubeconfig
	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig, "Path to the kubeconfig file")
	flag.Parse()

	fmt.Printf("Using kubeconfig: %s\n", kubeconfig)

	// Measure time for client creation
	clientStartTime := time.Now()

	// Build the config from the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	clientDuration := time.Since(clientStartTime)
	fmt.Printf("Time to create K8s client: %v\n", clientDuration)

	// Measure time for listing namespaces
	listStartTime := time.Now()

	// List namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing namespaces: %v\n", err)
		os.Exit(1)
	}

	listDuration := time.Since(listStartTime)
	fmt.Printf("Time to list namespaces: %v\n", listDuration)

	fmt.Println("Available namespaces:")
	for i, ns := range namespaces.Items {
		fmt.Printf("%d. %s\n", i+1, ns.Name)
	}
}
