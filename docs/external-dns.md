# External DNS

External DNS is a Kubernetes addon that configures public DNS servers with information about exposed Kubernetes services to make them discoverable. It allows you to control DNS records dynamically via Kubernetes resources.

We leverage external-dns to automatically create DNS records for our tenant API endpoints. This allows us to provide a seamless experience for our tenants to access their API endpoints.

## Prerequisites

- An Azure DNS Zone
- A Service Principal with the `DNS Zone Contributor` role on the Azure DNS Zone

## Installation

Create a json file with the contents below and save it as `external-dns-azure-secret.json`.

```json
{
"tenantId": "your-tenant-id",
"subscriptionId": "your-subscription-id",
"aadClientId": "your-client-id",
"aadClientSecret": "your-client-secret",
"resourceGroup": "resource group name where the DNS zone is located",
}
```

Create a secret in the `external-dns` namespace using the json file created above.

```bash
kubectl create ns external-dns
kubectl create secret generic azure-config-file --from-file=azure.json=external-dns-azure-secret.json -n external-dns
```

You'll also need a `txt-owner-id` configmap to be created in the `external-dns` namespace. This is used to identify the records created by external-dns. You can create this configmap by running the command below.

```bash
kubectl create configmap -n external-dns txt-owner-id --from-literal=txt-owner-id=unique-id-for-your-cluster
```
