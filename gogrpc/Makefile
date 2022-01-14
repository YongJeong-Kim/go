# Windows
gen:
	protoc --proto_path=proto --go_out=. proto/*.proto

clean:
	del pb\*.go

client:
	go run cmd/client/main.go

server:
	go run cmd/server/main.go -port 8080

# Linux
#gen:
#	protoc --proto_path=proto --go_out=plugins=grpc:pb proto/*.proto
#
#clean:
#	rm pb/*.go
