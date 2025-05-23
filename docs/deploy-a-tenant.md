# Deploy a tenant vcluster

```bash
git clone https://github.com/kubespaces-io/kubespaces-public.git
cd kubespaces-public/charts/vcluster

```bash
export tenant=meetup
export org=suse
export location_short=ne
export cloud=azure
export domain=kubespaces.cloud

helm upgrade -i -n $tenant-$org --create-namespace \
 $tenant-$org oci://ghcr.io/kubespaces-io/kubespaces-public/tenant --version 0.1.6 \
--set tenant.location_short=$location_short \
--set tenant.cloud=$cloud \
--set tenant.name=$tenant \
--set tenant.org=$org \
--set tenant.domain=$domain \
--set controlPlane.distro.k8s.version="v1.32.1" \
--set "vcluster.exportKubeConfig.server=https://api.$tenant.$org.$location_short.$cloud.$domain" \
--set "vcluster.controlPlane.proxy.extraSANs[0]=api.$tenant.$org.$location_short.$cloud.$domain"
```

kg po -n $tenant-$org -w

Access the tenant:

```bash
kubectl get secret vc-$tenant-$org -n $tenant-$org -o yaml -o jsonpath='{.data.config}' | base64 -d > /tmp/$tenant-$org-kubeconfig
export KUBECONFIG=/tmp/$tenant-$org-kubeconfig
kubectl get pods -A
kubectl get nodes
```