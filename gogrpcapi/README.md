## gRPC API
### [download gRPC](https://grpc.io/docs/languages/go/quickstart/#prerequisites])


```shell
$ protoc --proto_path=proto --go_out=. proto/*.proto
$ protoc --proto_path=<PROTO_DIRECTORY_PATH> --go_out=<OUT_DIRECTORY_PATH> <PROTO_FILE_PATH>
```

### [download gRPC-go](https://github.com/grpc/grpc-go#installation)
```shell
$ go get -u google.golang.org/grpc
```



### [install protocol buffer](https://github.com/protocolbuffers/protobuf/releases?page=1)

https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/adding_annotations/#using-protoc


### [installation grpc gateway](https://github.com/grpc-ecosystem/grpc-gateway#installation)
```shell
$ go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

### [github buf](https://github.com/bufbuild/buf)
- [docs buf](https://buf.build/docs/installation#windows-support)


### generate gRPC stubs ```buf generate```
```yaml
version: v1
plugins:
  - plugin: go
    out: gen/go
    opt:
      - paths=source_relative
  - plugin: go-grpc
    out: gen/go
    opt:
      - paths=source_relative
```

```shell
protoc -I . \
    --go_out ./gen/go/ --go_opt paths=source_relative \
    --go-grpc_out ./gen/go/ --go-grpc_opt paths=source_relative \
    your/service/v1/your_service.proto
```

```buf generate <PROTO_PATH>```
