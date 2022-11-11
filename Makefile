all:
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
