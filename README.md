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

- Go 1.24 or later
- Access to a Kubernetes cluster
- kubectl configured with access to your cluster

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/k8s-api-bench.git
   cd k8s-api-bench
   ```

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

Specify a namespace to use for the benchmarking operations:

```bash
./k8s-api-bench --namespace=kube-system
```

Combine both options:

```bash
./k8s-api-bench --kubeconfig=/path/to/your/kubeconfig --namespace=kube-system
```
