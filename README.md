[![CodeFactor](https://www.codefactor.io/repository/github/natron-io/tenant-api/badge)](https://www.codefactor.io/repository/github/natron-io/tenant-api)
# tenant-api
API to present data to the tenant-dashboard with a GitHub oauth login.
Tenants represents the teams of a GitHub organization.

## api
`/github/login` - Login with GitHub \
`/api/v1/pods` - Get pods of a tenant \
`/api/v1/namespaces` - Get namespaces of a tenant \
`/api/v1/serviceaccounts` - Get serviceaccounts of a tenant by namespaces \
`/api/v1/requests/cpu` - Get cpurequests in **Milicores** of a tenant \
`/api/v1/requests/memory` - Get memoryrequests in **Bytes** of a tenant \
`/api/v1/requests/storage` - Get storagerequests in **Bytes** of a tenant by storageclass \
`/api/v1/costs/cpu` - Get the cpu costs by  \

## env
`CLIENT_ID` - GitHub client id **required** \
`CLIENT_SECRET` - GitHub client secret **required** \
`SECRET_KEY` - JWT secret key *optional* (default: random 32 bytes, displayed in the logs) \
`LABELSELECTOR` - label key for selecting tenant ressources *optional* (default: "natron.io/tenant") \
`CALLBACK_URL` - GitHub oauth callback url without path *optional* (default: "http://localhost:3000") \
`CPU_COST` - Cost of a cpu in your currency *optional* (default: 1.00 for 1 CPU) \
`MEMORY_COST` - Cost of a memory in your currency *optional* (default: 1.00 for 1 GB) \
`STORAGE_COST_<storageclass name>` - Cost of your storage classes in your currency *optional* (default: 1.00 for 1 GB)

## local testing

```bash
minikube start
kubectl apply -f sa.yaml
kubectl apply -f rbac.yaml
kubectl apply -f deployment.yaml
kubectl expose deployment tenant-api --type=NodePort --port=8000

minikube service tenant-api
```