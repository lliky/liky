gen:
	protoc --proto_path=proto proto/*.proto  --go_out=pb --go-grpc_out=pb

clean:
	rm pb/*.go

server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go
