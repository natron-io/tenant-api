[![CodeFactor](https://www.codefactor.io/repository/github/natron-io/tenant-api/badge)](https://www.codefactor.io/repository/github/natron-io/tenant-api)
# tenant-api
API to present data to the tenant-dashboard with a GitHub oauth login.
Tenants represents the teams of a GitHub organization.

## api

### auth
`/github/login` - Login with GitHub \
`/github/login/callback` - Callback after GitHub login \
`/logout` - Logout \

### tenant ressource requests
`/api/v1/pods` - Get pods of a tenant \
`/api/v1/namespaces` - Get namespaces of a tenant \
`/api/v1/serviceaccounts` - Get serviceaccounts of a tenant by namespaces \
`/api/v1/requests/cpu` - Get cpurequests in **Milicores** of a tenant \
`/api/v1/requests/memory` - Get memoryrequests in **Bytes** of a tenant \
`/api/v1/requests/storage` - Get storagerequests in **Bytes** of a tenant by storageclass \

### tenant ressource costs
`/api/v1/costs/cpu` - Get the cpu costs by CPU \
`/api/v1/costs/memory` - Get the memory costs by Memory \
`/api/v1/costs/storage` - Get the storage costs by StorageClass \

## env

### GitHub
`CLIENT_ID` - GitHub client id **required** \
`CLIENT_SECRET` - GitHub client secret **required** \
`CALLBACK_URL` - GitHub oauth callback url without path *optional* (default: "http://localhost:3000") \

### auth
`SECRET_KEY` - JWT secret key *optional* (default: random 32 bytes, displayed in the logs) \

### tenant ressource identifiers
`LABELSELECTOR` - label key for selecting tenant ressources *optional* (default: "natron.io/tenant") \

### cost calculation values
`CPU_COST` - Cost of a cpu in your currency *optional* (default: 1.00 for 1 CPU) \
`MEMORY_COST` - Cost of a memory in your currency *optional* (default: 1.00 for 1 GB) \
`STORAGE_COST_<storageclass name>` - Cost of your storage classes in your currency *optional, multiple allowed* (default: 1.00 for 1 GB)

## local testing
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