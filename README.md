[![CodeFactor](https://www.codefactor.io/repository/github/natron-io/tenant-api/badge)](https://www.codefactor.io/repository/github/natron-io/tenant-api)
# tenant-api
API to present data to the tenant-dashboard with a GitHub oauth login.
Tenants represents the teams of a GitHub organization.

## api
`/github/login` - Login with GitHub \
`/api/v1/pods` - Get pods of a tenant \
`/api/v1/namespaces` - Get namespaces of a tenant \
`/api/v1/serviceaccounts` - Get serviceaccounts of a tenant by namespaces \
`/api/v1/cpurequests` - Get cpurequests of a tenant \
`/api/v1/memoryrequests` - Get memoryrequests of a tenant \
`/api/v1/storagerequests` - Get storagerequests of a tenant by storageclass \

## env
`CLIENT_ID` - GitHub client id **required** \
`CLIENT_SECRET` - GitHub client secret **required** \
`SECRET_KEY` - JWT secret key *optional* (default: random 32 bytes, displayed in the logs) \
`LABELSELECTOR` - label key for selecting tenant ressources *optional* (default: "natron.io/tenant")

## local testing

```bash
minikube start
kubectl apply -f sa.yaml
kubectl apply -f rbac.yaml
kubectl apply -f deployment.yaml
kubectl expose deployment tenant-api --type=NodePort --port=8000

minikube service tenant-api
```