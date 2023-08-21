.PHONY: dev build start

dev:
	go run ./*.go

build:
	go build  -o ./build/finetrack ./cmd/finetrack/main.go

start:
	./build/finetrack

proto:
	protoc --proto_path=protos --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	protos/*.proto
