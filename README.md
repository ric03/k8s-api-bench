# k8s-api-bench

A simple tool for benchmarking Kubernetes API operations.

## Description

k8s-api-bench is a Go-based utility that connects to a Kubernetes cluster and performs operations to help measure and
analyze the performance of the Kubernetes API server.

Currently, the tool can:

- Connect to a Kubernetes cluster using a kubeconfig file
- List available namespaces
- Measure the performance of various API operations used for tab completion in kubectl:
    - Listing pods in a namespace
    - Listing deployments in a namespace
    - Listing services in a namespace
    - Listing ConfigMaps in a namespace
    - Listing Secrets in a namespace
    - Listing API resources
    - Listing Custom Resource Definitions (simulated)

## Installation

### Prerequisites

- Go 1.22 or later (Note: The go.mod file currently specifies Go 1.24, which may need to be updated)
- Access to a Kubernetes cluster
- kubectl configured with access to your cluster

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/[your-username]/k8s-api-bench.git
   cd k8s-api-bench
   ```

   Note: Replace `[your-username]` with your actual GitHub username.

2. Build the binary:

   Using Go directly:
   ```bash
   go build -o k8s-api-bench
   ```

   Or using the provided Makefile:
   ```bash
   make build           # Build for current platform
   make build-all       # Build for all supported platforms (Linux and Windows, 64-bit)
   make build-linux-amd64   # Build for Linux 64-bit
   make build-windows-amd64 # Build for Windows 64-bit
   ```

   Additional Makefile targets:
   ```bash
   make clean           # Remove build directory
   make run             # Run the application
   make help            # Show all available targets
   ```

   All build artifacts will be placed in the `build/` directory.

## Usage

Run the tool with the default kubeconfig location and default namespace:

```bash
./k8s-api-bench
```

Specify a custom kubeconfig file:

```bash
./k8s-api-bench --kubeconfig=/path/to/your/kubeconfig
```

Specify the number of iterations for each benchmark operation:

```bash
./k8s-api-bench --iterations=10
```

Combine both options:

```bash
./k8s-api-bench --kubeconfig=/path/to/your/kubeconfig --iterations=10
```

Note: The tool automatically runs benchmarks on all available namespaces in the cluster.

## Example Output

The test is performed with a local kind cluster.

Below are the performance statistics for various Kubernetes API operations. Note that your actual results will vary depending on your cluster configuration, resources, and current load:

```
--- Benchmark Statistics ---
Operation                          |          Min |          Max |          Avg |       Median |          P95
-----------------------------------+--------------+--------------+--------------+--------------+--------------
list API resources                 |       2.0 ms |       4.3 ms |       2.8 ms |       2.7 ms |       4.3 ms
list ConfigMaps                    |     198.3 ms |     201.3 ms |     200.0 ms |     200.0 ms |     201.0 ms
list Custom Resource Definitions   |       1.3 ms |       1.7 ms |       1.4 ms |       1.4 ms |       1.7 ms
list Secrets                       |     198.7 ms |     201.3 ms |     200.0 ms |     200.0 ms |     200.9 ms
list all API resources             |       2.0 ms |       3.7 ms |       2.4 ms |       2.2 ms |       3.7 ms
list deployments                   |       1.3 ms |       2.5 ms |       1.7 ms |       1.7 ms |       2.3 ms
list namespaces                    |       1.3 ms |     183.5 ms |      19.6 ms |       1.3 ms |     183.5 ms
list pods                          |     191.5 ms |     207.9 ms |     200.0 ms |     200.0 ms |     201.1 ms
list services                      |     179.2 ms |     201.3 ms |     198.3 ms |     199.9 ms |     200.7 ms
```

These statistics show the performance characteristics of different API operations, including minimum, maximum, average,
median, and 95th percentile response times.
