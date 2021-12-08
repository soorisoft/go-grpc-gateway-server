package main

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"go-grpc-gateway-server/pkg/api"
)

type server struct {
	log *zap.SugaredLogger
	api.UnimplementedTestServicesServer
}


func (s *server) SayHello(ctx context.Context, request *api.SayHelloRequest) (*api.SayHelloResponse, error) {

	s.log.Info("Hit to SayHello service",zap.Any("request", request))
	response := &api.SayHelloResponse{Msg: "Hello " + request.Name + "!!"}
	fmt.Println("Returning resp to SayHello service, Response=", response.Msg)
	return response, nil
}

func (s *server) SayHi(ctx context.Context, request *api.SayHiRequest) (*api.SayHiResponse, error) {

	s.log.Info("Hit to Say Hi service", zap.Any("request", request))
	response := &api.SayHiResponse{Msg: "Hi " + request.Name + "!!"}
	s.log.Info("Returning resp to Say Hi service, Response=", response.Msg)
	return response, nil
}
