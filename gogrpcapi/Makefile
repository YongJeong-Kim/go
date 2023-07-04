# Windows
gen:
	protoc --proto_path=proto --go_out=pb --go-grpc_out=pb --go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	proto/*.proto

clean:
	del pb\*.go

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

publish:
	go run cmd/client/publish.go -address 0.0.0.0:8080

server:
	go run cmd/server/main.go -port 8080

test:
	go test -cover ./...

# Linux
#gen:
#	protoc --proto_path=proto --go_out=plugins=grpc:pb proto/*.proto
#
#clean:
#	rm pb/*.go
