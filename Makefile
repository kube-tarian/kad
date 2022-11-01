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
	docker compose -f ./docker-compose-mysql-es.yml up -d --no-recreate
	sleep 20
	go test -timeout 120s -run ^TestIntegration* github.com/kube-tarian/kad/deployment-worker/integration_tests -v

stop-docker-compose-test:
	docker compose -f ./docker-compose-mysql-es.yml down -v

build:
	go mod download
	CGO_ENABLED=0 go build -o build/deployment_worker deployment-worker/main.go

clean:
	rm -rf build

docker-build:
	docker build -f dockerfiles/deployment-worker/Dockerfile -t ${DEPLOYMENT_WORKER_APP_NAME}:${BUILD} .

start-manual-test:
	docker compose -f ./docker-compose-mysql-es.yml up -d --no-recreate

stop-manual-test:
	docker compose -f ./docker-compose-mysql-es.yml down -v
