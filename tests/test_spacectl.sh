#!/bin/bash

# Configuration variables
TENANT_NAME="xtenant"
TENANT_ORG="xtest-org"
TENANT_CLOUD="kind"
TENANT_DOMAIN="localhost"
TENANT_OUTPUT_FILE="test_tenant.yaml"
NS="${TENANT_NAME}-${TENANT_ORG}"
SVC="${TENANT_NAME}-${TENANT_ORG}"
VCLUSTER_PORT=8443
VCLUSTER_SERVICE_PORT=443

# Timeouts
VCLUSTER_TIMEOUT="120s"
PORT_FORWARD_DELAY=3

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
