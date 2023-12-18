# GO Auth

## Tool installation

The following tools are required to be installed to compile and execute the project

### Environment setup Go

```shell
export PATH=$PATH:$(go env GOPATH)/bin
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

### Prerequisites

Before we start coding, we have to install some tools.

We will be using a Go gRPC server in the examples, so please install Go first
from [https://golang.org/dl/](https://golang.org/dl/)

After installing Go, use `go install to download and build the following binaries:

```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

This installs the `protoc` generator plugins we need to generate the stubs. Make sure to add `$GOPATH/bin` to your `$PATH`
so that executables installed via `go get` are available on your `$PATH`.

## Google Apis Proto

- [google/api/annotations.proto](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto)
- [google/api/http.proto](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto)
