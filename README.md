# go-grpc-gateway-server

This repository provides an example for go based microservice. Go micro services developed based
on gRPC protobuf's and also uses gRPC gateway as a reverse proxy to support JSON clients/requests.

## Dev pre-requisites
1. golang
    ```
   brew install go
   Note: Make sure to have go1.17.2 or latest
    ```
2. Docker
    ```
   brew install docker
   ```
   or [downlooad docker from here](https://docs.docker.com/get-docker/)


3. PgAdmin for Postgres DB, [download](https://www.pgadmin.org/download/) and this is optional if you prefer cli

4. Update env
    ```
   vi ~/.zshrc

    add below entries, save & exit.

    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    export GONOPROXY="github.com,gopkg.in,go.uber.org"
    ```
    ```
    source ~/.zshrc
    ```

## Install protoc

If you have Homebrew (which you can get from https://brew.sh), just run:

``brew install protobuf``

If you see any error messages, run brew doctor, follow any recommended fixes, and try again. If it still fails, try instead:

``brew upgrade protobuf``

Alternately, run the following commands:

```
PROTOC_ZIP=protoc-3.14.0-osx-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP```
```

```
go install \
github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
google.golang.org/protobuf/cmd/protoc-gen-go@latest \
google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Note: Make sure protobuf's installed correctly by running ```protoc --version```

## Generate protobuf's 
Note: Required this step only when you make changes to .proto file or introduce a new one
```
cd pkg/api

go generate ./...

(or) alternative try below in case you did not set GOPATH

env PATH="$HOME/go/bin:$PATH" go generate ./...
```

Following 3 files are auto generated ones with above command
```
api/test_api.pb.go
api/test_api.pb.gw.go
api/test_api_grpc.pb.go
```

## Build and Run Docker

```bash
docker compose build

docker compose up
```

## Instructions for Developers

### How to add a new service & test?
This section demonstrates, how to add a new grpc gateway service, implement and test it in local.

Note: Follow naming standards from here: https://cloud.google.com/apis/design/naming_convention#method_names

####1. Add new proto changes
Lets add a test api called `SayHello` to `pkg/api/test_api.proto`

```
// SayHello service generates a greeting message
rpc SayHello(SayHelloRequest) returns (SayHelloResponse) {
 option (google.api.http) = {
   post: "/v1/greet"
   body: "*"
 };
}
```
   Add request `SayHelloRequest` & response `SayHelloResponse` objects declared above

   ```
 // Request Message for greeting service
 message SayHelloRequest {
   string name = 1;
 }

// Response Message for greeting service
message SayHelloResponse {
  string msg = 1;
}
```

####2. generate protobufs
```
cd pkg/api

go generate ./...

(or) alternative try below in case you did not set GOPATH

env PATH="$HOME/go/bin:$PATH"
```

####3. Implement service handler
```
func (s *Server) SayHello(ctx context.Context, request *api.SayHelloRequest) (*api.SayHelloResponse, error) {

	s.Log.Info("Hit to SayHello service and Request=%s", request.Name)
	response := &api.SayHelloResponse{Msg: "Hello " + request.Name + "!!"}
	return response, nil
}
```

####4. Build, Deploy and Test
```
docker compose build

docker compose up
```

```
curl -X POST -d '{"name": "John Doe" }' -H 'Content-Type: application/json' 
http://localhost:9091/v1/greet
```

Response:
```
{"msg": "Hello John Doe!!"}
```
