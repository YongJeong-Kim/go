# Windows
#gen:
#	protoc --proto_path=proto \
#	--go_out=pb --go_opt=paths=source_relative \
#	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
#	--grpc-gateway_out=allow_delete_body=true:pb --grpc-gateway_opt=paths=source_relative \
#	proto/*.proto

gen:
	buf generate proto

clear:
	rmdir /S /Q pb

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
