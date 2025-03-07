# Connect your tenant to AAD

export org=bldrnet
export location_short=ne
export cloud=azure
export domain=kubespaces.cloud

export tenant=projz

vcluster:
  controlPlane:
    distro:
      k8s:
        apiServer:
          enabled: true
          # ExtraArgs are additional arguments to pass to the distro binary.
          extraArgs:
          - --oidc-client-id=eaa9c642-6353-49e3-88b1-45735ce82c51
          - --oidc-issuer-url=https://sts.windows.net/30c10907-19b2-4fb8-9272-a4a539628560/
          - --oidc-username-claim=email
          - --oidc-groups-claim=groups

helm upgrade -i -n $tenant-$org --create-namespace \
$tenant-$org . \
--set tenant.location_short=$location_short \
--set tenant.cloud=$cloud \
--set tenant.name=$tenant \
--set tenant.org=$org \
--set tenant.domain=$domain \
--set controlPlane.distro.k8s.version="v1.31.1" \
--set "vcluster.exportKubeConfig.server=https://api.$tenant.$org.$location_short.$cloud.$domain" \
--set "vcluster.controlPlane.proxy.extraSANs[0]=api.$tenant.$org.$location_short.$cloud.$domain" --values ~/Desktop/temp.aad.yaml

kubectl wait po -n $tenant-$org $tenant-$org-0 --for=condition=ready

kubectl get secret vc-$tenant-$org -n $tenant-$org -o json | jq -r '.data.config'| base64 -D > ~/temp/vclusters/$tenant-$org.yaml
export KUBECONFIG=~/temp/vclusters/$tenant-$org.yaml

kubectl apply -f - <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: test-sso-clusterrolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: "7a24fa70-9cd9-4284-bf45-12af68e1327b"
EOF


cp ~/temp/vclusters/$tenant-$org.yaml ~/temp/vclusters/$tenant-$org-aad.yaml
sed -i '' 's/kubernetes-admin/azuread/g' ~/temp/vclusters/${tenant}-${org}-aad.yaml
sed -i '' '/client-certificate-data/,$c\
    exec:\
      apiVersion: client.authentication.k8s.io/v1beta1\
      args:\
      - oidc-login\
      - get-token\
      - --oidc-issuer-url=https://sts.windows.net/30c10907-19b2-4fb8-9272-a4a539628560/\
      - --oidc-client-id=eaa9c642-6353-49e3-88b1-45735ce82c51\
      - --oidc-client-secret=vYV8Q~eDKZAqBGZsv3rS7kH.aLkqOtaBScEFscle\
      command: kubectl\
      env: null\
      interactiveMode: IfAvailable\
      provideClusterInfo: false'  ~/temp/vclusters/${tenant}-${org}-aad.yaml

export KUBECONFIG=~/temp/vclusters/$tenant-$org-aad.yaml

