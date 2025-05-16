#!/bin/bash

# Configuration variables
CLUSTER_NAME="kubespaces-host"
GATEWAY_API_VERSION="v1.0.0"
CERT_MANAGER_VERSION="v1.14.0"
ISTIO_VERSION="1.26.0"
ISTIO_DIR="istio-${ISTIO_VERSION}"
ISTIO_PROFILE="ambient"
TENANT_NAME="xtenant"
TENANT_ORG="xtest-org"
TENANT_CLOUD="kind"
TENANT_DOMAIN="localhost"
TENANT_OUTPUT_FILE="test_tenant.yaml"
NS="${TENANT_NAME}-${TENANT_ORG}"
SVC="${TENANT_NAME}-${TENANT_ORG}"
VCLUSTER_PORT=8443
VCLUSTER_SERVICE_PORT=443

# Gateway API URLs
GATEWAY_API_STANDARD_URL="https://github.com/kubernetes-sigs/gateway-api/releases/download/${GATEWAY_API_VERSION}/standard-install.yaml"
GATEWAY_API_EXPERIMENTAL_URL="https://github.com/kubernetes-sigs/gateway-api/releases/download/${GATEWAY_API_VERSION}/experimental-install.yaml"
CERT_MANAGER_URL="https://github.com/cert-manager/cert-manager/releases/download/${CERT_MANAGER_VERSION}/cert-manager.yaml"

# Timeouts
CERT_MANAGER_TIMEOUT="60s"
VCLUSTER_TIMEOUT="120s"
PORT_FORWARD_DELAY=3

# Create Kind cluster if it doesn't exist
if ! kind get clusters | grep -q "${CLUSTER_NAME}"; then
  echo "Creating Kind cluster..."
  kind create cluster --name "${CLUSTER_NAME}" --config kind-config.yaml
  sleep 5
fi

# Install Gateway API components
echo "Installing Gateway API components..."
kubectl apply -f "${GATEWAY_API_STANDARD_URL}"
kubectl apply -f "${GATEWAY_API_EXPERIMENTAL_URL}"

# Install cert-manager
echo "Installing cert-manager..."
kubectl apply -f "${CERT_MANAGER_URL}"
kubectl -n cert-manager wait --for=condition=ready pod -l app=cert-manager --timeout="${CERT_MANAGER_TIMEOUT}"

# Create necessary namespaces
echo "Creating required namespaces..."
kubectl create ns istio-system

# Install Istio
echo "Installing Istio..."
if [ ! -d "${ISTIO_DIR}" ]; then
  echo "Downloading Istio ${ISTIO_VERSION}..."
  curl -L https://istio.io/downloadIstio | ISTIO_VERSION="${ISTIO_VERSION}" sh -
fi

./"${ISTIO_DIR}"/bin/istioctl install --set profile="${ISTIO_PROFILE}" -y

# Deploy tenant using spacectl
echo "Deploying test tenant using spacectl..."
../spacectl/spacectl tenant create "${TENANT_NAME}" \
  -l local \
  --cloud "${TENANT_CLOUD}" \
  --tenant "${TENANT_NAME}" \
  --org "${TENANT_ORG}" \
  --domain "${TENANT_DOMAIN}" \
  --output-file "${TENANT_OUTPUT_FILE}" \
  --wait

# Wait for vcluster to be ready
echo "Waiting for vcluster to be ready..."
kubectl wait -n "${NS}" --for=condition=ready pod -l app=vcluster --timeout="${VCLUSTER_TIMEOUT}"

# Start port-forwarding in the background
echo "Starting port-forwarding..."
kubectl -n "${NS}" port-forward svc/"${SVC}" "${VCLUSTER_PORT}:${VCLUSTER_SERVICE_PORT}" &
PF_PID=$!

# Give port-forwarding time to start
sleep "${PORT_FORWARD_DELAY}"

# Print instructions
echo ""
echo "vCluster is ready! To use it:"
echo "export KUBECONFIG=$(pwd)/${TENANT_OUTPUT_FILE}"
echo ""
echo "When done, kill the port-forwarding process with:"
echo "kill ${PF_PID}"
