#!/bin/bash
# Destroy kind cluster
kind delete cluster --name kubespaces-host
# Remove tenant kubeconfig
rm ./test_tenant.yaml
# Remove istio folder 
rm -rf ./istio-*