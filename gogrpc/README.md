## gRPC
### [download gRPC](https://grpc.io/docs/languages/go/quickstart/#prerequisites])


```shell
$ protoc --proto_path=proto --go_out=. proto/*.proto
$ protoc --proto_path=<PROTO_DIRECTORY_PATH> --go_out=<OUT_DIRECTORY_PATH> <PROTO_FILE_PATH>
```