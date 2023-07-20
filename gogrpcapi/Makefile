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

mock:
	mockgen -package accountv1mock -destination pb/account/v1/mock/account_service_grpc_mock.go github.com/YongJeong-Kim/go/gogrpcapi/pb/account/v1 AccountServiceServer

test:
	go test -cover ./...



# Linux
#gen:
#	protoc --proto_path=proto --go_out=plugins=grpc:pb proto/*.proto
#
#clean:
#	rm pb/*.go
