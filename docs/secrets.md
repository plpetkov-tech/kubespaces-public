# Create secrets for external-dns, cert-manager, and other components

kubectl create ns external-dns
kubectl create secret generic azure-config-file -n external-dns --from-file=azure.json=/Users/alessandro/repos/scratch/kubeconfigs/external-dns.json
kubectl create configmap -n external-dns txt-owner-id --from-literal=txt-owner-id=susemeetup
kubectl create ns cert-manager
export CLIENT_SECRET="secret" #from az ad sp create-for-rbac --name extdns --role "DNS Zone Contributor" --scopes /subscriptions/1c51d1c3-d83d-4d71-ace1-df3496eddac4/resourceGroups/dns -o json

kubectl create secret generic azuredns-config -n cert-manager --from-literal=client-secret=${CLIENT_SECRET}