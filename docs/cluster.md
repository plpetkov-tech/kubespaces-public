# Create a cluster for Kubespaces

## Azure

```bash
export CLUSTER_RG=kubespaces
export CLUSTER_NAME=kubespaces
export AKS_VERSION=1.32.0
export BASE_NODE_COUNT=2
export NODE_MAX=5
export NODE_MIN=2
export LOCATION=northeurope
export NODE_SIZE=Standard_D16s_v6
export TENANT_NODE_MIN_COUNT=0
export TENANT_NODE_MAX_COUNT=10
export TENANT_NODE_SIZE=Standard_D16s_v6

```

```bash
az group create --name $CLUSTER_RG --location $LOCATION
```

```bash
az aks create \
    --resource-group $CLUSTER_RG \
    --name $CLUSTER_NAME \
    --node-count $BASE_NODE_COUNT \
    --min-count 1 \
    --max-count 5 \
    --node-vm-size $NODE_SIZE \
    --generate-ssh-keys \
    --location $LOCATION \
    --enable-cluster-autoscaler \
    --kubernetes-version $AKS_VERSION \
    --auto-upgrade-channel rapid \
    --enable-aad \
    --enable-image-cleaner \
    --network-plugin azure \
    --network-plugin-mode overlay \
    --pod-cidr 192.168.0.0/16 \
    --network-dataplane cilium \
    --node-os-upgrade-channel NodeImage \
    --os-sku AzureLinux \
    --enable-managed-identity \
    --api-server-authorized-ip-ranges 95.99.46.198/32 \
    --enable-addons azure-policy \
    --enable-defender 

az aks nodepool add \
    -g $CLUSTER_RG \
    --cluster-name $CLUSTER_NAME \
    --node-count 0 \
    --enable-cluster-autoscaler \
    --min-count $TENANT_NODE_MIN_COUNT \
    --max-count $TENANT_NODE_MAX_COUNT \
    --priority Spot \
    --os-type Linux  \
    --os-sku AzureLinux \
    --node-vm-size $TENANT_NODE_SIZE \
    --labels "priority=spot","workload=tenants"  \
    -n tenants \
    --mode User \
    --vm-set-type VirtualMachineScaleSets --no-wait
```

Get the credentials to the cluster:

```bash
az aks get-credentials --resource-group $CLUSTER_RG --name $CLUSTER_NAME --admin
```

Install flux (needs the flux cli installed):

```bash
flux install
```

Create the secrets before continuing, check secrets.md for more information.

Add the flux repo and the first kustomization:

```bash
kubectl apply -f gitops/host/host-gitrepo.yaml
kubectl apply -f gitops/host/host-kustomization.yaml
```

This deploys everything needed for the host cluster.

```bash
kubectl get kustomizations.kustomize.toolkit.fluxcd.io -A
kubectl get helmreleases.kustomize.toolkit.fluxcd.io -A
```

```bash
kubectl get gateway -n istio-system
```

The LoadBalancer IP is the public IP of the Istio Gateway; it will be attached to the DNS record for the tenant API.

Note:

Due to this [issue](https://github.com/kubernetes-sigs/azuredisk-csi-driver/issues/2777), you need to run this command:

```bash
kubectl apply -f https://raw.githubusercontent.com/andyzhangx/demo/refs/heads/master/aks/download-v6-disk-rules.yaml
```
