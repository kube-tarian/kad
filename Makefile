PREFIX := kad
SERVER_APP_NAME := server
AGENT_APP_NAME := agent
DEPLOYMENT_WORKER_APP_NAME := deployment-worker
CONFIG_WORKER_APP_NAME := config-worker
BUILD := 0.1.1

gen-protoc:
	cd proto && protoc --go_out=../agent/pkg/agentpb/ --go_opt=paths=source_relative \
    		--go-grpc_out=../agent/pkg/agentpb/ --go-grpc_opt=paths=source_relative \
    		./agent.proto

	cd proto && protoc --go_out=../server/pkg/pb/agentpb --go_opt=paths=source_relative \
    		--go-grpc_out=../server/pkg/pb/agentpb --go-grpc_opt=paths=source_relative \
    		./agent.proto
	
	cd proto && protoc --go_out=../server/pkg/pb/climonpb --go_opt=paths=source_relative \
    		--go-grpc_out=../server/pkg/pb/climonpb --go-grpc_opt=paths=source_relative \
    		./climon.proto

	cd proto && protoc --go_out=../climon/pkg/pb/climonpb --go_opt=paths=source_relative \
                --go-grpc_out=../climon/pkg/pb/climonpb --go-grpc_opt=paths=source_relative \
                ./climon.proto

docker-build-server:
	# The prefix for server to changed either as server or intelops-kad-server
	docker build -f dockerfiles/server/Dockerfile -t ${PREFIX}-${SERVER_APP_NAME}:${BUILD} .

docker-build-kad: docker-build-agent docker-build-deployment docker-build-config

docker-build-agent:
	docker build -f dockerfiles/agent/Dockerfile -t ${PREFIX}-${AGENT_APP_NAME}:${BUILD} .

docker-build-deployment:
	docker build -f dockerfiles/deployment-worker/Dockerfile -t ${PREFIX}-${DEPLOYMENT_WORKER_APP_NAME}:${BUILD} .

docker-build-config:
	docker build -f dockerfiles/config-worker/Dockerfile -t ${PREFIX}-${CONFIG_WORKER_APP_NAME}:${BUILD} .

docker-build: docker-build-kad docker-build-server
