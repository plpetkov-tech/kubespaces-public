# Deploy a tenant vcluster

```bash
cd charts/vcluster

```bash
export tenant=bldrcon
export org=bldrnet
export location_short=ne
export cloud=azure
export domain=kubespaces.cloud

helm upgrade -i -n $tenant-$org --create-namespace \
$tenant-$org . \
--set tenant.location_short=$location_short \
--set tenant.cloud=$cloud \
--set tenant.name=$tenant \
--set tenant.org=$org \
--set tenant.domain=$domain \
--set controlPlane.distro.k8s.version="v1.31.1" \
--set "vcluster.exportKubeConfig.server=https://api.$tenant.$org.$location_short.$cloud.$domain" \
--set "vcluster.controlPlane.proxy.extraSANs[0]=api.$tenant.$org.$location_short.$cloud.$domain"
```

Access the tenant:

```bash
kubectl get secret $tenant-$org-kubeconfig -n $tenant-$org -o jsonpath='{.data.kubeconfig}' | base64 -d > /tmp/$tenant-$org-kubeconfig
export KUBECONFIG=/tmp/$tenant-$org-kubeconfig
kubectl get pods
```