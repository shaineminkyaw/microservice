proto:  
	# rm -rf pb/*.go
	protoc --proto_path=./user/user_proto --go_out=./pb/ --go_opt=paths=source_relative \
    --go-grpc_out=./pb/ --go-grpc_opt=paths=source_relative \
    ./user/user_proto/*.proto
server:
	go run ./cmd/server/grpc_server.go

evans:
	evans --host localhost --port 9099 -r repl

cert:
	./cert/gen.sh

.PHONY: proto server evans cert