# Demo flow for kubespaces

- Deploy an host cluster
- Deploy Flux and its components


kubectl run -n default nginx --image=nginx --port=80; kubectl expose po/nginx
kubectl apply -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: nginx-route
  namespace: default  # Adjust if your nginx service is in a different namespace
  annotations:
    external-dns.alpha.kubernetes.io/hostname: nginx-tenant.apps.ne.azure.kubespaces.cloud
spec:
  hostnames:
    - nginx-tenant.apps.ne.azure.kubespaces.cloud
  parentRefs:
    - name: apps
      namespace: istio-system
      sectionName: apps  # This maps to the listener name
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /
      backendRefs:
        - name: nginx-x-default-x-suse-meetup
          port: 80
EOF


** Bugs:

- [ ] This particularly nasty bug with Azure LoadBalancer: https://github.com/Azure/AKS/issues/3646
- 