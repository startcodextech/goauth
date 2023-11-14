# GO Auth

## Tool installation
The following tools are required to be installed to compile and execute the project

### Environment setup Go
```shell
export PATH=$PATH:$(go env GOPATH)/bin
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

### protoc
You can install the protocol compiler, `protoc`, with a package manager under Linux or macOS using the following commands.

- Linux, using `apt` or `apt-get`, for example:
```shell
apt install -y protobuf-compiler
```
- MacOS, using [Homebrew](https://brew.sh/)
```shell
brew install protobuf
```

### swag

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

### protoc-gen-openapi
```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
```
### protoc-gen-grpc-gateway
```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
```