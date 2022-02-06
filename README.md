<p align="center">
  <img src="docs/images/tenant-api-screenshot.png" />
</p>

# tenant-api
[![CodeFactor](https://www.codefactor.io/repository/github/natron-io/tenant-api/badge)](https://www.codefactor.io/repository/github/natron-io/tenant-api)
![Build Status](https://github.com/natron-io/tenant-api/workflows/CI/badge.svg) 
[![Go Report Card](https://goreportcard.com/badge/github.com/natron-io/tenant-api)](https://goreportcard.com/report/github.com/natron-io/tenant-api) 
![GitHub top language](https://img.shields.io/github/languages/top/natron-io/tenant-api)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/natron-io/tenant-api) 
![open issues](https://img.shields.io/github/issues-raw/natron-io/tenant-api)
![license](https://img.shields.io/github/license/natron-io/tenant-api)

API to present data to the [tenant-dashboard](https://github.com/natron-io/tenant-dashboard) with a GitHub oauth login.
Tenants represents the teams of a GitHub organization.

## api

#### `GET`
> **important:** for authenticated access you need to provide the `Authorization` header with the `Bearer` token.

You can add `<tenant>` in front of the path to get the tenant specific data (of everything). 
> e.g. `/api/v1/<tenant>/pods`
#### auth
`/login/github` - Login with GitHub \
`/login/github/callback` - Callback after GitHub login

#### notifications
`/api/v1/notifications` - Get the Slack notification messages of the broadcast channel provided via envs

##### general tenant resources
`/api/v1/<tenant>/pods` - Get pods of a tenant \
`/api/v1/<tenant>/namespaces` - Get namespaces of a tenant \
`/api/v1/<tenant>/serviceaccounts` - Get serviceaccounts of a tenant by namespaces \

##### specific tenant resources
`/api/v1/<tenant>/requests/cpu` - Get cpurequests in **Milicores** of a tenant \
`/api/v1/<tenant>/requests/memory` - Get memoryrequests in **Bytes** of a tenant \
`/api/v1/<tenant>/requests/storage` - Get storagerequests in **Bytes** of a tenant by storageclass \
`/api/v1/<tenant>/requests/ingress` - Get ingress resourcess total of a tenant by ingressclass

##### tenant resources costs
`/api/v1/<tenant>/costs/cpu` - Get the CPU costs by CPU \
`/api/v1/<tenant>/costs/memory` - Get the memory costs by Memory \
`/api/v1/<tenant>/costs/storage` - Get the storage costs by StorageClass \
`/api/v1/<tenant>/costs/ingress` - Get the ingress costs by tenant

##### tenant resource quotas
`/api/v1/<tenant>/quotas/cpu` - Get the CPU resource Quota by the label defined via env \
`/api/v1/<tenant>/quotas/memory` - Get the memory resource Quota by the label defined via env \
`/api/v1/<tenant>/quotas/storage` - Get the storage resource Quota for each storage class by the label**s** defined via env 


#### `POST`

##### auth
You can send the github code with json body `{"github_code": "..."}` to the `/login/github` endpoint.
> The code you need to generate must have the `read:org` scope.

## env

### general
`CORS` - CORS middleware for Fiber that that can be used to enable Cross-Origin Resource Sharing with various options. (e.g. "https://example.com, https://example2.com")

### GitHub
> There are two ways for authenticating with GitHub. You can authenticate without a dashboard, so the github callback url is not the same as the dashboard.

`CLIENT_ID` - GitHub client id **required** \
`CLIENT_SECRET` - GitHub client secret **required** \
`CALLBACK_URL` - GitHub oauth callback url without path *optional* (default: "http://localhost:3000")

### auth
`SECRET_KEY` - JWT secret key *optional* (default: random 32 bytes, displayed in the logs)

### notifications
`SLACK_TOKEN` - Tenant API Slack Application User Token *optional* (if not set, the notification REST route will be deactivated) \
`SLACK_BROADCAST_CHANNEL_ID` - BroadCast Slack Channel ID *optional* (**required** if SLACK_TOKEN is set)

### tenant resources identifiers
`TENANT_LABEL` - label key for selecting tenant resourcess *optional* (default: "natron.io/tenant")

### cost calculation values
`DISCOUNT_LABEL` - label key for selecting the discount value *optional* (default: "natron.io/discount" (float -> e.g. "0.1")) \
`CPU_COST` - Cost of a CPU in your currency *optional* (default: 1.00 for 1 CPU) \
`MEMORY_COST` - Cost of a memory in your currency *optional* (default: 1.00 for 1 GB) \
`STORAGE_COST_<storageclass name>` - Cost of your storage classes in your currency **required, multiple allowed** (default: 1.00 for 1 GB) \
`INGRESS_COST` - Cost of ingress in your currency *optional* (default: 1.00 for 1 ingress)

### resource quotas
`QUOTA_NAMESPACE_SUFFIX` - The namespace suffix where the config of the tenant configuration takes place *optional* (default: "config" e.g. namespace name: "tenant-config") \
`QUOTA_CPU_LABEL` - The CPU quota label *optional* (default: "natron.io/cpu-quota") \
`QUOTA_MEMORY_LABEL` - The memory quota label *optional* (default: "natron.io/memory-quota")
`QUOTA_STORAGE_LABEl_<storageclass name>` - The storage label of each storage class *optional, multiple allowed* (default: "natron.io/storage-quota-<storageclass name>" renders every storageclass defined in the Storage)

## labels

### resource quotas
For setting the quota labels on the tenant config namespace, you have to enter the values in the following format:  
- CPU: `cores` e.g. natron.io/cpu-quota: "1" (-> 1 Core)
- Memory: `GB` e.g. natron.io/memory-quota: "4" (-> 4GB)
- Storage: `GB` e.g. natron.io/storage-quota-<storageclass name>: "50" (-> 50GB)

## deployment
*example deployment files:* [kubernetes manifests](docs/kubernetes)

1. run a local minikube and apply a service account with clusterwide `view` permissions
```bash
minikube start
kubectl apply -f sa.yaml
kubectl apply -f rbac.yaml
kubectl apply -f deployment.yaml
kubectl expose deployment tenant-api --type=NodePort --port=8000

minikube service tenant-api
```
2. create a GitHub application in your GitHub organization and set the url (and port displayed at exposing the service via minikube) to the `CALLBACK_URL` (without path) and for the callback URL set the `CALLBACK_URL` with the path `/login/github/callback` (e.g. http://localhost:3000/login/github/callback)