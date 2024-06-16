all: proto user voting

proto:
	@protoc --go_out=buffers/ --go_opt=paths=source_relative \
    --go-grpc_out=buffers/ --go-grpc_opt=paths=source_relative \
    protobuffs/user.proto

user:
	# @go build -o bin/user ./cmd/user
	@go build -o user ./cmd/user

voting:
	# @go build -o bin/user ./cmd/voting
	@go build -o voting ./cmd/voting