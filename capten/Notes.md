# Pre-requisite

```
# Start kind kubernetes environment
# Assumption: Base folder path go/github.com/kube-tarian/kad
$ cd integrator/
$ make start-kind

# Download and install temporal
# Reference: https://github.com/temporalio/helm-charts

$ cd ~/go/src/github.com/temporalio/helm-charts
$ helm install     --set server.replicaCount=1     --set cassandra.config.cluster_size=1     --set prometheus.enabled=false     --set grafana.enabled=false     --set elasticsearch.enabled=false     temporal . --timeout 8m
```

## Expose temporal

```
$ kubectl port-forward service/temporal-web 9082:8080
```

# Install server component in kind

```
# Create custom server values file /tmp/server-values.yaml
$ vim /tmp/server-values.yaml
image:
  repository: kad-server
  pullPolicy: IfNotPresent
  tag: "0.1.1"

cassandra:
  host: temporal-cassandra.default.svc.cluster.local

$ cd charts/server/
$ helm install server . -f /tmp/server-values.yaml
```

## Expose server component

```
$ kubectl port-forward services/server -n default 9092:8080
```

# Install kad component in kind

```
# Create custom kad values file /tmp/kad-values.yaml
$ vim /tmp/kad-values.yaml

image:
  pullPolicy: IfNotPresent
  tag: "0.1.1"
agent:
  repository: kad-agent
deployment_worker:
  repository: kad-deployment-worker
config_worker:
  repository: kad-config-worker

ingressroute:
  enabled: false

ingress:
  enabled: false

$ cd charts/kad/
$ helm install kad . -f /tmp/kad-values.yaml
```

# Install argo-cd using integrator framework helm plugin

```
# Open temporal web UI (http://127.0.0.1:9082) and check the workflow activity of this action
```

## Expose argocd details

```
kubectl port-forward services/argocd-server -n default 9081:443

# Set ARGOCD service url env variable
$ export ARGOCD_SERVICE_URL=https://localhost:9081

# Get argocd admin password
$ export PASSWORD=$(kubectl -n default get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)

```

# Install hello-world application using integrator framework argocd plugin

```
# make kind-integration-argocd-test

# Open temporal web UI (http://127.0.0.1:9082) and check the workflow activity of this action
```

# Delete argocd using deployer

```
# make kind-integration-helm-delete-test

# Open temporal web UI and check the workflow activity of this action
```

# Delete temporal component

```
$ helm delete temporal
```

# Delete kind (k8s) cluster

```
$ cd integrator/
$ make stop-kind
```
