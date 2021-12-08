package api

//go:generate protoc -I../../.. -I../../third_party/grpc-gateway-v2.3.0/third_party/googleapis --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=logtostderr=true test_api.proto
