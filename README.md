# Kubespaces

Kubespaces is a multi-tenant Kubernetes platform that enables secure and isolated tenant environments using virtual clusters. It provides a GitOps-based deployment model with comprehensive networking, security, and management capabilities.

![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)

## Overview

Kubespaces creates isolated tenant environments using virtual clusters (vclusters) with:

- Secure network isolation using Gateway API and Istio Ambient Mesh
- Automatic DNS and TLS certificate management
- Tenant-specific resource quotas and limits
- Simple CLI management tool (`spacectl`)
- GitOps-based deployment with Flux CD

## Architecture

The platform consists of a host Kubernetes cluster that manages multiple tenant virtual clusters:

- **Host Cluster**: Contains shared infrastructure components (Istio, Gateway API, cert-manager, external-DNS)
- **Tenant vClusters**: Isolated Kubernetes environments for each tenant with their own API server and control plane
- **Gateway API**: Provides controlled network access to tenant environments
- **Istio Ambient Mesh**: Enables secure service mesh capabilities

## Prerequisites

- Kubernetes cluster v1.28+ (AKS, EKS, GKE, or kind for local development)
- Helm 3.10+
- kubectl 1.28+
- Flux CD v2 (for GitOps deployments)
- A domain name with DNS zone management (for external-dns)

## Installation

### Setting up the Host Cluster

1. Clone the repository:
```bash
git clone https://github.com/kubespaces-io/kubespaces-public.git
cd kubespaces-public
```

2. For AKS deployment, follow the instructions in [docs/cluster.md](docs/cluster.md):
```bash
# Example for creating an AKS cluster:
export CLUSTER_RG=demo
export CLUSTER_NAME=kubespaces
export LOCATION=northeurope
# See docs/cluster.md for full instructions
```

3. Install Flux CD:
```bash
flux install
```

4. Create required secrets (see [docs/secrets.md](docs/secrets.md)):
```bash
# Example: Create secrets for external-dns and cert-manager
kubectl create ns external-dns
kubectl create secret generic azure-config-file -n external-dns --from-file=azure.json=/path/to/external-dns.json
kubectl create configmap -n external-dns txt-owner-id --from-literal=txt-owner-id=your-cluster-id
```

5. Deploy the GitOps components:
```bash
kubectl apply -f gitops/host/host-gitrepo.yaml
kubectl apply -f gitops/host/host-kustomization.yaml
```

## Usage

### Managing Tenants with spacectl

The `spacectl` CLI tool simplifies tenant management:

#### Creating a Tenant

```bash
# Create a tenant with default settings
./spacectl tenant create --tenant dev --org contoso

# Create a tenant with custom settings
./spacectl tenant create \
  --tenant dev \
  --org contoso \
  --cloud azure \
  --location-short ne \
  --domain kubespaces.cloud \
  --k8s-version 1.31.1 \
  --wait \
  --output-file kubeconfig.yaml
```

#### Updating a Tenant

```bash
./spacectl tenant update \
  --tenant dev \
  --org contoso \
  --k8s-version 1.32.0
```

#### Deleting a Tenant

```bash
./spacectl tenant delete --tenant dev --org contoso
```

### Accessing Tenant Clusters

After creating a tenant, you can access it using the generated kubeconfig:

```bash
# Export the kubeconfig
kubectl get secret vc-dev-contoso -n dev-contoso -o yaml -o jsonpath='{.data.config}' | base64 -d > dev-contoso-kubeconfig
export KUBECONFIG=dev-contoso-kubeconfig

# Verify access
kubectl get pods -A
```

## Testing

Run the included test scripts to verify your deployment:

```bash
# Test GitOps deployment with kind
./tests/test_gitops.sh

# Test spacectl functionality
./tests/test_spacectl.sh
```

## Documentation

Detailed documentation is available in the `docs/` directory:

- [Cluster Setup](docs/cluster.md) - Creating a Kubernetes cluster for Kubespaces
- [Deploy a Tenant](docs/deploy-a-tenant.md) - Step-by-step guide for tenant deployment
- [External DNS](docs/external-dns.md) - Setting up DNS for tenant access
- [AAD Integration](docs/aad.md) - Connecting tenants to Azure Active Directory

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](./LICENSE) file for details.

## ToDo

- [ ] Add more details to the README
- [ ] Use [patches](https://www.vcluster.com/docs/vcluster/0.20.0/configure/vcluster-yaml/experimental/generic-sync?x1=1#patches-reference) to rewrite the HTTPRoute in the tenant vcluster to a proper format in the host cluster
- [ ] Add a demo for the tenant vcluster
- [ ] Add the HBONE Istio Ambient port (15008) to the NetworkPolicy created bu vcluster helm chart
- [ ] Improve CLI safety checks and error messages
