# KubeSpaces Helm Chart

This repository contains the Helm chart for deploying KubeSpaces tenants. A tenant in KubeSpaces is implemented as a vCluster running on a host cluster, providing isolated Kubernetes environments for developers and platform engineers.

## Overview

The KubeSpaces Helm chart enables the deployment of tenant vClusters with the following features:
- Isolated Kubernetes environments using vCluster
- Gateway API integration for ingress management
- Network policies for tenant isolation
- RBAC configuration for namespace access control

## Prerequisites

- Kubernetes cluster (v1.24+)
- Helm 3.x
- Gateway API CRDs installed on the host cluster
- vCluster operator installed on the host cluster

## Installation

1. Add the KubeSpaces Helm repository:
```bash
helm repo add kubespaces https://kubespaces-io.github.io/kubespaces-public
helm repo update
```

2. Install a tenant:
```bash
helm install my-tenant kubespaces/tenant \
  --namespace my-tenant \
  --create-namespace \
  --set tenant.name=my-tenant \
  --set tenant.organization=my-org
```

## Configuration

Key configuration parameters:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `tenant.name` | Name of the tenant | `""` (required) |
| `tenant.organization` | Organization ID | `""` (required) |
| `tenant.resources.limits.cpu` | CPU limit for tenant workloads | `"4"` |
| `tenant.resources.limits.memory` | Memory limit for tenant workloads | `"8Gi"` |
| `tenant.networkPolicy.enabled` | Enable network policies | `true` |

For a complete list of configuration options, see [values.yaml](./helm/tenant/values.yaml).

## Usage

After installation, you can:
1. Get the tenant's kubeconfig:
```bash
kubectl get secret -n my-tenant my-tenant-kubeconfig -o jsonpath='{.data.config}' | base64 -d > kubeconfig.yaml
```

2. Use the kubeconfig to interact with your tenant:
```bash
kubectl --kubeconfig kubeconfig.yaml get namespaces
```

## Security

- Each tenant runs in an isolated vCluster
- Network policies restrict tenant-to-tenant communication
- RBAC roles are automatically created for namespace access control
- Tenant resources are limited by default quotas

## Testing

The repository includes a test suite that helps you verify the Helm chart functionality locally. The tests use Kind (Kubernetes in Docker) to create a local cluster and deploy a test tenant.

### Prerequisites for Testing

- Docker
- Kind
- kubectl
- Helm 3.x
- curl

### Running Tests

1. The test suite will:
   - Create a Kind cluster named `kubespaces-host`
   - Install required components (Gateway API, cert-manager, Istio)
   - Deploy a test tenant
   - Set up port-forwarding for local access

2. Run the test script:
```bash
cd tests
./test_kubespaces.sh
```

3. After the test completes, you can access the test tenant:
```bash
# Use the test tenant's kubeconfig
export KUBECONFIG=$(pwd)/test_tenant.yaml

# Verify the tenant is working
kubectl get namespaces
```

4. Clean up when done:
```bash
# Kill the port-forwarding process
kill $PF_PID  # The PID is shown in the test output

# Optional: Delete the Kind cluster
kind delete cluster --name kubespaces-host
```

### Test Environment Details

The test environment includes:
- A local Kind cluster with Gateway API support
- Istio ambient mesh for service mesh capabilities
- cert-manager for certificate management
- A test tenant with:
  - Isolated vCluster environment
  - Network policies
  - RBAC configuration
  - Local domain (localhost) setup

The test tenant is configured with minimal resources suitable for local testing. For production deployments, adjust the resource limits and other parameters as needed.

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](./LICENSE) file for details.

## Roadmap

- [ ] Implement HTTPRoute patches for proper host cluster integration
- [ ] Add Istio Ambient (HBONE) port (15008) to NetworkPolicy
- [ ] Create tenant deployment demo
- [ ] Improve CLI safety checks and error messages
