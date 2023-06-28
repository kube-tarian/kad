PREFIX := kad
SERVER_APP_NAME := server
AGENT_APP_NAME := agent
DEPLOYMENT_WORKER_APP_NAME := deployment-worker
CONFIG_WORKER_APP_NAME := config-worker
CLIMON_APP_NAME := climon
BUILD := 0.1.1

gen-protoc:
	mkdir -p server/pkg/pb/serverpb
	mkdir -p capten/agent/pkg/agentpb
	mkdir -p server/pkg/pb/agentpb
	mkdir -p server/pkg/pb/climonpb

	cd proto && protoc --go_out=../server/pkg/pb/serverpb/ --go_opt=paths=source_relative \
    		--go-grpc_out=../server/pkg/pb/serverpb  --go-grpc_opt=paths=source_relative \
    		./server.proto

	cd proto && protoc --go_out=../capten/agent/pkg/agentpb/ --go_opt=paths=source_relative \
    		--go-grpc_out=../capten/agent/pkg/agentpb/ --go-grpc_opt=paths=source_relative \
    		./agent.proto

	cd proto && protoc --go_out=../server/pkg/pb/agentpb --go_opt=paths=source_relative \
    		--go-grpc_out=../server/pkg/pb/agentpb --go-grpc_opt=paths=source_relative \
    		./agent.proto
	
	cd proto && protoc --go_out=../capten/capten-sdk/agentpb/ --go_opt=paths=source_relative \
    		--go-grpc_out=../capten/capten-sdk/agentpb/ --go-grpc_opt=paths=source_relative \
    		./agent.proto

	cd proto && protoc --go_out=../server/pkg/pb/climonpb --go_opt=paths=source_relative \
    		--go-grpc_out=../server/pkg/pb/climonpb --go-grpc_opt=paths=source_relative \
    		./climon.proto

docker-build-server:
	# The prefix for server to changed either as server or intelops-kad-server
	docker build --platform=linux/amd64 -f dockerfiles/server/Dockerfile -t ${PREFIX}-${SERVER_APP_NAME}:${BUILD} .

docker-build-kad: docker-build-agent docker-build-deployment docker-build-config

docker-build-agent:
	docker build -f dockerfiles/agent/Dockerfile -t ${PREFIX}-${AGENT_APP_NAME}:${BUILD} .

docker-build-deployment:
	docker build -f dockerfiles/deployment-worker/Dockerfile -t ${PREFIX}-${DEPLOYMENT_WORKER_APP_NAME}:${BUILD} .

docker-build-config:
	docker build -f dockerfiles/config-worker/Dockerfile -t ${PREFIX}-${CONFIG_WORKER_APP_NAME}:${BUILD} .

docker-build-climon:
	docker build -f dockerfiles/climon/Dockerfile -t ${PREFIX}-${CLIMON_APP_NAME}:${BUILD} .

docker-build: docker-build-kad docker-build-server docker-build-climon
