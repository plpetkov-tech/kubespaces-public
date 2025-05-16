#!/bin/bash

# Create Kind cluster if it doesn't exist
if ! kind get clusters | grep -q kubespaces-host; then
  echo "Creating Kind cluster..."
  kind create cluster --name kubespaces-host --config kind-config.yaml
fi
sleep 5
# Install Gateway API CRDs
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/standard-install.yaml
# Install Gateway API experimental CRDs
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/experimental-install.yaml
# Install External-DNS (optional for local testing, or use mock DNS)
# For local testing, we can use a ConfigMap to simulate DNS resolution

# Install cert-manager (for TLS certificates)
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.0/cert-manager.yaml
kubectl -n cert-manager wait --for=condition=ready pod -l app=cert-manager --timeout=60s


# Create necessary namespaces
kubectl create ns istio-system

# Install Istio (lightweight profile for local testing)
curl -L https://istio.io/downloadIstio | sh -
cd istio-*/bin
./istioctl install --set profile=ambient -y

# Deploy tenant
echo "Deploying test tenant..."
helm upgrade -i -n test-tenant-org --create-namespace \
  test-tenant ../charts/tenant \
  --set tenant.location_short=local \
  --set tenant.cloud=kind \
  --set tenant.name=test \
  --set tenant.org=org \
  --set tenant.domain=localhost \
  --set "vcluster.exportKubeConfig.server=https://localhost:8443" \
  --set "vcluster.controlPlane.proxy.extraSANs[0]=localhost"

# Wait for tenant to be ready
echo "Waiting for tenant to be ready..."
kubectl wait -n test-tenant-org --for=condition=ready pod -l app=vcluster --timeout=120s

# Get kubeconfig
echo "Getting kubeconfig..."
kubectl get secret -n test-tenant-org vc-test-tenant -o jsonpath='{.data.config}' | base64 -d > test_tenant.yaml

# Start port-forwarding in the background
echo "Starting port-forwarding..."
kubectl -n test-tenant-org port-forward svc/test-tenant 8443:443 &
PF_PID=$!

# Give port-forwarding time to start
sleep 3

# Print instructions
echo ""
echo "vCluster is ready! To use it:"
echo "export KUBECONFIG=$(pwd)/test_tenant.yaml"
echo ""
echo "When done, kill the port-forwarding process with:"
echo "kill $PF_PID"
