gen:
	protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/*.proto

clean:
	del pb\*.go

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080 -to ccc -from ddd