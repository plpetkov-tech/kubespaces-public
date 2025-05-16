#!/bin/bash

set -euo pipefail

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# Get the workspace root directory (parent of tests directory)
WORKSPACE_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print status messages
print_status() {
    echo -e "${GREEN}==>${NC} $1"
}

# Function to print warning messages
print_warning() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

# Function to print error messages
print_error() {
    echo -e "${RED}Error:${NC} $1"
}

# Function to wait for pods to be ready
wait_for_pods() {
    local namespace=$1
    local label=$2
    local timeout=$3
    print_status "Waiting for pods with label $label in namespace $namespace to be ready..."
    kubectl wait --for=condition=ready pod -l "$label" -n "$namespace" --timeout="$timeout" || {
        print_error "Pods with label $label in namespace $namespace failed to become ready within $timeout"
        kubectl get pods -n "$namespace" -l "$label"
        exit 1
    }
}

# Create Kind cluster if it doesn't exist
if ! kind get clusters | grep -q gitops-host; then
    print_status "Creating Kind cluster..."
    kind create cluster --name gitops-host --config "$SCRIPT_DIR/kind-config.yaml"
else
    print_warning "Kind cluster 'gitops-host' already exists"
fi

# Wait for cluster to be ready
sleep 5

# Create flux-system namespace
print_status "Installing Flux in flux-system namespace..."
helm install -n flux-system --create-namespace flux oci://ghcr.io/fluxcd-community/charts/flux2

# Apply infrastructure components in order
print_status "Applying Gitops Manifests..."
kubectl apply -f "$WORKSPACE_ROOT/gitops/host/host-gitrepo.yaml"
kubectl apply -f "$WORKSPACE_ROOT/gitops/host/host-kustomization.yaml"

# Verify the deployment
print_status "Verifying deployment..."
echo "Checking Gateway API resources:"
kubectl get gatewayclasses
kubectl get gateways -A

echo "Checking Istio resources:"
kubectl get pods -n istio-system

echo "Checking cert-manager resources:"
kubectl get pods -n cert-manager

echo "Checking external-dns resources:"
kubectl get pods -n external-dns

echo "Checking tenant resources:"
kubectl get tenants -A

print_status "GitOps test deployment completed successfully!"
echo ""
echo "To clean up the cluster, run:"
echo "kind delete cluster --name gitops-host"
echo ""
echo "To use the cluster:"
echo "kubectl config use-context kind-gitops-host" 
