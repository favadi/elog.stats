default:
	@mkdir -p bin/server
	@mkdir -p bin/client
	@go build -o bin/server/elog cmd/server/*.go
	@go build -o bin/client/elog cmd/client/*.go
	@echo Done
grpc:
	@protoc -I proto proto/elog.proto --go_out=plugins=grpc:./
	@echo Done