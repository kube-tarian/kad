all:
	cd proto && protoc --go_out=../cilmon/pkg/pb/climonpb --go_opt=paths=source_relative \
    		--go-grpc_out=../cilmon/pkg/pb/climonpb --go-grpc_opt=paths=source_relative \
    		./climon.proto
