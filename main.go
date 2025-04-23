package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/discovery"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// BenchmarkResults stores the results of all benchmark operations
type BenchmarkResults struct {
	// Map of operation name to slice of durations
	Results map[string][]time.Duration
}

// NewBenchmarkResults creates a new BenchmarkResults instance
func NewBenchmarkResults() *BenchmarkResults {
	return &BenchmarkResults{
		Results: make(map[string][]time.Duration),
	}
}

// Add adds a new duration for the specified operation
func (br *BenchmarkResults) Add(operation string, duration time.Duration) {
	br.Results[operation] = append(br.Results[operation], duration)
}

// Helper function to measure the execution time of a function
func measureTime(name string, f func() error, results *BenchmarkResults) {
	startTime := time.Now()
	err := f()
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("Error during %s: %v\n", name, err)
	} else {
		fmt.Printf("Time to %s: %v\n", name, duration)
		// Store the duration in the results
		results.Add(name, duration)
	}
}

// Helper function to run a benchmark operation multiple times
func runBenchmark(name string, iterations int, f func() error, results *BenchmarkResults) {
	fmt.Printf("Running benchmark '%s' for %d iterations...\n", name, iterations)
	for i := 0; i < iterations; i++ {
		fmt.Printf("Iteration %d/%d: ", i+1, iterations)
		measureTime(name, f, results)
	}
}

// Calculate statistics for the benchmark results
func (br *BenchmarkResults) CalculateStats() map[string]map[string]time.Duration {
	stats := make(map[string]map[string]time.Duration)

	for op, durations := range br.Results {
		if len(durations) == 0 {
			continue
		}

		// Sort durations for percentile calculations
		sort.Slice(durations, func(i, j int) bool {
			return durations[i] < durations[j]
		})

		// Calculate statistics
		var sum time.Duration
		min := durations[0]
		max := durations[0]

		for _, d := range durations {
			sum += d
			if d < min {
				min = d
			}
			if d > max {
				max = d
			}
		}

		avg := sum / time.Duration(len(durations))

		// Calculate median (50th percentile)
		median := durations[len(durations)/2]
		if len(durations)%2 == 0 {
			median = (durations[len(durations)/2-1] + durations[len(durations)/2]) / 2
		}

		// Calculate 95th percentile
		p95Index := int(math.Ceil(float64(len(durations))*0.95)) - 1
		if p95Index >= len(durations) {
			p95Index = len(durations) - 1
		}
		p95 := durations[p95Index]

		// Store statistics
		stats[op] = map[string]time.Duration{
			"min":    min,
			"max":    max,
			"avg":    avg,
			"median": median,
			"p95":    p95,
		}
	}

	return stats
}

// formatDuration formats a time.Duration to show only one decimal place in milliseconds
func formatDuration(d time.Duration) string {
	// Convert to milliseconds with one decimal place
	ms := float64(d.Microseconds()) / 1e3
	return fmt.Sprintf("%.1f ms", ms)
}

// Print the statistics in a readable format
func (br *BenchmarkResults) PrintStats() {
	stats := br.CalculateStats()

	// Sort operations for consistent output
	operations := make([]string, 0, len(stats))
	for op := range stats {
		operations = append(operations, op)
	}
	sort.Strings(operations)

	// Calculate the maximum length of operation names
	maxOpLength := 0
	for _, op := range operations {
		if len(op) > maxOpLength {
			maxOpLength = len(op)
		}
	}

	// Add some padding to the maximum length
	opColWidth := maxOpLength + 2

	// Define column width for time values
	timeColWidth := 12

	fmt.Println("\n--- Benchmark Statistics ---")

	// Create the header with dynamic width
	headerFormat := fmt.Sprintf("%%-%ds | %%%ds | %%%ds | %%%ds | %%%ds | %%%ds\n",
		opColWidth, timeColWidth, timeColWidth, timeColWidth, timeColWidth, timeColWidth)
	fmt.Printf(headerFormat, "Operation", "Min", "Max", "Avg", "Median", "P95")

	// Create the separator line with dynamic width
	separatorLine := strings.Repeat("-", opColWidth) + "-+" +
		strings.Repeat("-", timeColWidth+2) + "+" +
		strings.Repeat("-", timeColWidth+2) + "+" +
		strings.Repeat("-", timeColWidth+2) + "+" +
		strings.Repeat("-", timeColWidth+2) + "+" +
		strings.Repeat("-", timeColWidth+2)
	fmt.Println(separatorLine)

	// Create the row format with dynamic width
	rowFormat := fmt.Sprintf("%%-%ds | %%%ds | %%%ds | %%%ds | %%%ds | %%%ds\n",
		opColWidth, timeColWidth, timeColWidth, timeColWidth, timeColWidth, timeColWidth)

	for _, op := range operations {
		stat := stats[op]
		fmt.Printf(rowFormat,
			op,
			formatDuration(stat["min"]),
			formatDuration(stat["max"]),
			formatDuration(stat["avg"]),
			formatDuration(stat["median"]),
			formatDuration(stat["p95"]))
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
	var iterations int

	// If the kubeconfig flag is not provided, use the default path
	home := homedir.HomeDir()
	if home == "" {
		fmt.Println("Error: unable to find home directory")
		os.Exit(1)
	}
	defaultKubeconfig := filepath.Join(home, ".kube", "config")

	// Set up the flags
	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig, "Path to the kubeconfig file")
	flag.IntVar(&iterations, "iterations", 1, "Number of iterations for each benchmark operation")
	flag.Parse()

	if iterations < 1 {
		fmt.Println("Error: iterations must be at least 1")
		os.Exit(1)
	}

	fmt.Printf("Using kubeconfig: %s\n", kubeconfig)
	fmt.Printf("Running each benchmark operation for %d iterations\n", iterations)

	// Create benchmark results object
	benchmarkResults := NewBenchmarkResults()

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

	// Get namespaces (we need this for later operations)
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing namespaces: %v\n", err)
		os.Exit(1)
	}

	// Benchmark listing namespaces
	runBenchmark("list namespaces", iterations, func() error {
		_, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		return err
	}, benchmarkResults)

	fmt.Println("Available namespaces:")
	for i, ns := range namespaces.Items {
		fmt.Printf("%d. %s\n", i+1, ns.Name)
	}

	// Benchmark operations used for tab completion
	fmt.Println("\n--- Tab Completion API Operations Benchmark ---")

	// Perform namespace-specific operations for each namespace
	for _, ns := range namespaces.Items {
		nsName := ns.Name
		fmt.Printf("\n--- Benchmarking namespace: %s ---\n", nsName)

		// List pods in the current namespace
		runBenchmark(fmt.Sprintf("list pods in namespace %s", nsName), iterations, func() error {
			return listPods(clientset, nsName)
		}, benchmarkResults)

		// List deployments in the current namespace
		runBenchmark(fmt.Sprintf("list deployments in namespace %s", nsName), iterations, func() error {
			return listDeployments(clientset, nsName)
		}, benchmarkResults)

		// List services in the current namespace
		runBenchmark(fmt.Sprintf("list services in namespace %s", nsName), iterations, func() error {
			return listServices(clientset, nsName)
		}, benchmarkResults)

		// List ConfigMaps in the current namespace
		runBenchmark(fmt.Sprintf("list ConfigMaps in namespace %s", nsName), iterations, func() error {
			return listConfigMaps(clientset, nsName)
		}, benchmarkResults)

		// List Secrets in the current namespace
		runBenchmark(fmt.Sprintf("list Secrets in namespace %s", nsName), iterations, func() error {
			return listSecrets(clientset, nsName)
		}, benchmarkResults)
	}

	// Non-namespace specific operations
	fmt.Println("\n--- Non-namespace specific operations ---")

	// List API resources
	runBenchmark("list API resources", iterations, func() error {
		return listAPIResources(clientset)
	}, benchmarkResults)

	// List all API resources
	runBenchmark("list all API resources", iterations, func() error {
		return listAllAPIResources(clientset)
	}, benchmarkResults)

	// List Custom Resource Definitions
	runBenchmark("list Custom Resource Definitions", iterations, func() error {
		return listCRDs(config)
	}, benchmarkResults)

	fmt.Println("\nBenchmarking complete!")

	// Print the benchmark statistics
	benchmarkResults.PrintStats()
}
