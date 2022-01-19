gen:
	protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/*.proto

clean:
	del pb\*.go

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080 -to ccc -from ddd

subsa:
	go run cmd/client/main.go -address 0.0.0.0:8080 -id aaa

subsb:
	go run cmd/client/main.go -address 0.0.0.0:8080 -id bbb

subsc:
	go run cmd/client/main.go -address 0.0.0.0:8080 -id ccc

sendafromb:
	go run cmd/client/main.go -address 0.0.0.0:8080 -to aaa -from bbb -content hello I am bbb

sendafromc:
	go run cmd/client/main.go -address 0.0.0.0:8080 -to aaa -from ccc -content hello I am ccc

sendbfromc:
	go run cmd/client/main.go -address 0.0.0.0:8080 -to bbb -from ccc -content hello I am ccc

sendallfroma:
	go run cmd/client/main.go -address 0.0.0.0:8080 -from aaa -content hello everyone I am aaa