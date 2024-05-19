# Overview

## Integrator components

- Integrator
    - deployment-worker
    - config-worker
- Agent (gRPC)

## IntelOps SaaS components

- Server

## IntelOps Market Place

- Temporal
- Postgres

# Installation using kind

## Prepare docker images

```
make docker-build
```

## Setup kind

### create kind

```
# Create a file kind-dev-cluster-config.yaml with below content, where "192.168.73.3" is Ip of the machine kind k8s to be created:

kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: "192.168.73.3"
  apiServerPort: 8443

# Start kind with the above configuration
kind create cluster --name dev-cluster --config kind-dev-cluster-config.yaml
```

### Delete kind

```
kind delete --name dev-cluster
```

## Load the images to kind

```
kind load docker-image kad-deployment-worker:0.1.1  --name dev-cluster
kind load docker-image kad-config-worker:0.1.1  --name dev-cluster
kind load docker-image kad-agent:0.1.1  --name dev-cluster
kind load docker-image kad-server:0.1.1  --name dev-cluster
```

## Setup server

### Install server

```
helm install server ./charts/server/
```

### Delete server

```
helm delete server
```

## Setup kad

### Install kad

```
helm install kad ./charts/kad/
```

### Delete kad

```
helm delete kad
```

## Setup temporal

### Install temporal

```
# Clone temporal helm charts
git clone https://github.com/temporalio/helm-charts.git

# Update dependancies
helm dependencies update

# Install temporal

### Delete temporal

```
helm delete temporal
```

## Portward service ports

```
kubectl port-forward services/temporal-frontend 7233:7233
kubectl port-forward service/temporal-web 9082:8080
kubectl port-forward services/server 9091:8080
```

## Testing helm plugin

### Update kubeconfiguration in below files

- "integrator/tests/argocd-helm-plugin.json"
- "integrator/tests/argocd-helm-plugin-delete.json"

### Create argocd app

```
curl -X POST http://127.0.0.1:9091/deploy -H "content-type: application/json" -d @integrator/tests/argocd-helm-plugin.json
```

### Delete argocd app

```
curl -X POST http://127.0.0.1:9091/deploy -H "content-type: application/json" -d @integrator/tests/argocd-helm-plugin-delete.json
```

## Testing argocd plugin

### Port forward the aprgocd for GUI access

```
kubectl port-forward svc/argocd-server  8083:80
```

### Get the argocd password

```
kubectl -n default get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

### Add below env parameters to kad-deployment-worker deployment according to argocd deployment

```
kubectl edit deployments.apps kad-deployment-worker

            - name: ARGOCD_SERVICE_URL
              value: {{ .Values.argocd.serviceURL }}
            - name: ARGOCD_PASSWORD
              value: {{ .Values.argocd.password }}
```

### Create helloworld app

```
curl -X POST http://127.0.0.1:9091/deploy -H "content-type: application/json" -d @tests/hello-world-argocd-plugin.json
```

### Delete helloworld app

```
curl -X POST http://127.0.0.1:9091/deploy -H "content-type: application/json" -d @tests/hello-world-argocd-plugin-delete.json

```

