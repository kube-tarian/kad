# Stage-1 Demo content

## Pre-requisite

```
# Start kind kubernetes environment
make start-kind

# Install temporal
cd /home/share2kanna/go/src/github.com/temporalio/helm-charts
helm install     --set server.replicaCount=1     --set cassandra.config.cluster_size=1     --set prometheus.enabled=false     --set grafana.enabled=false     --set elasticsearch.enabled=false     temporal . --timeout 8m 
```

## Expose temporal

```
kubectl port-forward services/temporal-frontend 7233:7233
kubectl port-forward service/temporal-web 9082:8080
```

## Install argo-cd using integrator framework helm plugin

```
# make kind-integration-helm-test
curl -X POST http://127.0.0.1:9092/deploy -H "content-type: application/json" -d @tests/argocd-helm-plugin.json 

# Open temporal web UI and check the workflow activity of this action
```

## Expose argocd details

```
kubectl port-forward services/argocd-server -n default 9081:443

# Set ARGOCD service url env variable
export ARGOCD_SERVICE_URL=https://localhost:9081

# Get argocd admin password
export PASSWORD=$(kubectl -n default get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)

```

## Install hello-world application using integrator framework argocd plugin

```
# make kind-integration-argocd-test

curl -X POST http://127.0.0.1:9092/deploy -H "content-type: application/json" -d @tests/hello-world-argocd-plugin.json 

# Open temporal web UI and check the workflow activity of this action
```

## Delete hello-world using argocd plugin

```
# make kind-integration-argocd-delete-test

curl -X POST http://127.0.0.1:9092/deploy -H "content-type: application/json" -d @tests/hello-world-argocd-plugin-delete.json

# Open temporal web UI and check the workflow activity of this action
```

## Delete argocd using helm plugin

```
# make kind-integration-helm-delete-test

curl -X POST http://127.0.0.1:9092/deploy -H "content-type: application/json" -d @tests/argocd-helm-plugin-delete.json 

# Open temporal web UI and check the workflow activity of this action
```

# integrator OEM details and its exposed ports

## Deployer worker listens on 9080 port

## Install temporal

```
cd /home/share2kanna/go/src/github.com/temporalio/helm-charts

helm install     --set server.replicaCount=1     --set cassandra.config.cluster_size=1     --set prometheus.enabled=false     --set grafana.enabled=false     --set elasticsearch.enabled=false     temporal . --timeout 8m 
```

## Expose temporal

```
kubectl port-forward services/temporal-frontend 7233:7233
```

## Expose temporal UI

```
kubectl port-forward service/temporal-web 9082:8080
```

## Expose arog-cd server UI

```
cd /home/share2kanna/go/src/github.com/argoproj/argo-helm
helm install argocd .
kubectl port-forward services/argocd-server -n default 9081:443
```
curl -X 'POST''localhost:9092/register/agent'-H 'accept: */*' -H 'Content-Type: application/json'{"customer_id": "1","endpoint":"kad-agent.default.svc.cluster.local:8080"}'