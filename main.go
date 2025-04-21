package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/discovery"
	"os"
	"path/filepath"
	"time"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Helper function to measure the execution time of a function
func measureTime(name string, f func() error) {
	startTime := time.Now()
	err := f()
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("Error during %s: %v\n", name, err)
	} else {
		fmt.Printf("Time to %s: %v\n", name, duration)
	}
}

// List pods in a namespace (used for tab completion)
func listPods(clientset *kubernetes.Clientset, namespace string) error {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Found %d pods in namespace %s\n", len(pods.Items), namespace)
	return nil
}

// List deployments in a namespace (used for tab completion)
func listDeployments(clientset *kubernetes.Clientset, namespace string) error {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Found %d deployments in namespace %s\n", len(deployments.Items), namespace)
	return nil
}

// List services in a namespace (used for tab completion)
func listServices(clientset *kubernetes.Clientset, namespace string) error {
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Found %d services in namespace %s\n", len(services.Items), namespace)
	return nil
}

// List ConfigMaps in a namespace (used for tab completion)
func listConfigMaps(clientset *kubernetes.Clientset, namespace string) error {
	configMaps, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Found %d ConfigMaps in namespace %s\n", len(configMaps.Items), namespace)
	return nil
}

// List Secrets in a namespace (used for tab completion)
func listSecrets(clientset *kubernetes.Clientset, namespace string) error {
	secrets, err := clientset.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Found %d Secrets in namespace %s\n", len(secrets.Items), namespace)
	return nil
}

// List API resources (used for tab completion)
func listAPIResources(clientset *kubernetes.Clientset) error {
	apiResources, err := clientset.Discovery().ServerPreferredResources()
	if err != nil {
		return err
	}

	resourceCount := 0
	for _, list := range apiResources {
		resourceCount += len(list.APIResources)
	}

	fmt.Printf("Found %d API resources\n", resourceCount)
	return nil
}

// List all API resources (used for tab completion)
func listAllAPIResources(clientset *kubernetes.Clientset) error {
	_, apiResources, err := clientset.Discovery().ServerGroupsAndResources()
	if err != nil {
		// Ignore group discovery errors, which happen when a resource isn't fully defined
		if !discovery.IsGroupDiscoveryFailedError(err) {
			return err
		}
		fmt.Printf("Warning: Some groups couldn't be discovered: %v\n", err)
	}

	resourceCount := 0
	for _, list := range apiResources {
		resourceCount += len(list.APIResources)
	}

	fmt.Printf("Found %d API resources (all)\n", resourceCount)
	return nil
}

// List Custom Resource Definitions (used for tab completion)
func listCRDs(config *rest.Config) error {
	// Create the apiextensions clientset
	apiextensionsClient, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating apiextensions client: %v", err)
	}

	// List CRDs
	crds, err := apiextensionsClient.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing CRDs: %v", err)
	}

	fmt.Printf("Found %d Custom Resource Definitions\n", len(crds.Items))
	return nil
}

func main() {
	// Define command-line flags
	var kubeconfig string
	var namespace string

	// If the kubeconfig flag is not provided, use the default path
	home := homedir.HomeDir()
	if home == "" {
		fmt.Println("Error: unable to find home directory")
		os.Exit(1)
	}
	defaultKubeconfig := filepath.Join(home, ".kube", "config")

	// Set up the flags
	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig, "Path to the kubeconfig file")
	flag.StringVar(&namespace, "namespace", "default", "Namespace to use for operations")
	flag.Parse()

	fmt.Printf("Using kubeconfig: %s\n", kubeconfig)
	fmt.Printf("Using namespace: %s\n", namespace)

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

	// Benchmark operations used for tab completion
	fmt.Println("\n--- Tab Completion API Operations Benchmark ---")

	// List pods in the specified namespace
	measureTime(fmt.Sprintf("list pods in namespace %s", namespace), func() error {
		return listPods(clientset, namespace)
	})

	// List deployments in the specified namespace
	measureTime(fmt.Sprintf("list deployments in namespace %s", namespace), func() error {
		return listDeployments(clientset, namespace)
	})

	// List services in the specified namespace
	measureTime(fmt.Sprintf("list services in namespace %s", namespace), func() error {
		return listServices(clientset, namespace)
	})

	// List API resources
	measureTime("list API resources", func() error {
		return listAPIResources(clientset)
	})

	// List all API resources
	measureTime("list all API resources", func() error {
		return listAllAPIResources(clientset)
	})

	// List ConfigMaps in the specified namespace
	measureTime(fmt.Sprintf("list ConfigMaps in namespace %s", namespace), func() error {
		return listConfigMaps(clientset, namespace)
	})

	// List Secrets in the specified namespace
	measureTime(fmt.Sprintf("list Secrets in namespace %s", namespace), func() error {
		return listSecrets(clientset, namespace)
	})

	// List Custom Resource Definitions
	measureTime("list Custom Resource Definitions", func() error {
		return listCRDs(config)
	})

	fmt.Println("\nBenchmarking complete!")
}
