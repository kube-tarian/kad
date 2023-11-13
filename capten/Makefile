DEPLOYMENT_WORKER_APP_NAME := deployment-worker
BUILD := 0.1.1

OPEN_API_CODEGEN := github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

${OPEN_API_CODEGEN}:
	$(eval TOOL=$(@:%=%))
	@echo Installing ${TOOL}...
	go install $(@:%=%)

tools: ${OPEN_API_CODEGEN}

OPEN_API_DIR = ./api

oapi-gen: tools oapi-gen-deployment-worker

oapi-gen-deployment-worker:
	$(eval APP_NAME=deployment-worker)
	@echo Generating server for ${APP_NAME}
	@mkdir -p ${APP_NAME}/${OPEN_API_DIR}
	${GOBIN}/oapi-codegen -config ./${APP_NAME}/cfg.yaml ./${APP_NAME}/openapi.yaml

start-docker-compose-test:
	docker compose -f ./docker-compose-cass-es.yml up -d --no-recreate
	sleep 20
	go test -coverprofile=cover.out -timeout 120s -run ^TestIntegration* github.com/kube-tarian/kad/integrator/deployment-worker/integration_tests -v

stop-docker-compose-test:
	docker compose -f ./docker-compose-cass-es.yml down -v

build:
	go mod download
	CGO_ENABLED=0 go build -o build/deployment_worker deployment-worker/main.go
	CGO_ENABLED=0 go build -o build/config_worker config-worker/main.go
	CGO_ENABLED=0 go build -o build/agent agent/cmd/main.go

clean:
	rm -rf build

start-kind:
	kind create cluster --config=tests/kind-config.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	# helm install     --set server.replicaCount=1     --set cassandra.config.cluster_size=1     --set prometheus.enabled=false     --set grafana.enabled=false     --set elasticsearch.enabled=false     temporal . --timeout 8m

stop-kind:
	kind delete cluster

kind-integration-argocd-test:
	go test -coverpkg=./... -coverprofile=cover-argocd.out -timeout 120s -run ^TestIntegrationArgocdDepl* ./... -v

kind-integration-argocd-delete-test:
	go test -coverpkg=./... -coverprofile=cover-argocd-delete.out -timeout 120s -run ^TestIntegrationArgocdDelete* ./... -v

kind-integration-helm-test:
	go test -coverpkg=./... -coverprofile=cover-helm.out -timeout 120s -run ^TestIntegrationHelmDepl* ./... -v

kind-integration-helm-delete-test:
	go test -coverpkg=./... -coverprofile=cover-helm-delete.out -timeout 120s -run ^TestIntegrationHelmDelete* ./... -v
