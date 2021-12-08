package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"go-grpc-gateway-server/pkg/api"
)

const (
	grpcPort   = ":50000"
	httpPort   = ":9091"
)

func main() {

	fx.New(
		fx.Provide(
			newLogger,
		),
		fx.Invoke(
			StartGrpcServer,
			StartGrpcGatewayServer1,
		),
	).Run()
}

// Returns basic version of logger from zap, refer: https://pkg.go.dev/go.uber.org/zap
func newLogger() *zap.SugaredLogger {
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()

	sugar.Info("Initialized new logger")
	return sugar
}

// StartGrpcServer starts the GRPC server
func StartGrpcServer(lc fx.Lifecycle, log *zap.SugaredLogger) {
	grpcServer := grpc.NewServer()
	s := &server{log: log}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			api.RegisterTestServicesServer(grpcServer, s)
			log.Info("Starting GRPC server, port", grpcPort)
			go func() {
				listener, err := net.Listen("tcp", grpcPort)
				if err != nil {
					log.Fatal("Unable to listen on TCP port,", zap.String("port", grpcPort), zap.Error(err))
				}
				if err = grpcServer.Serve(listener); err != nil {
					log.Fatal("GRPC server could not attach to listener.", zap.Error(err))
				}
			}()

			log.Info("============ GRPC server started.============ ")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			log.Info("============ GRPC server stopped.============ ")
			return nil
		},
	})
}

func StartGrpcGatewayServer1(lc fx.Lifecycle, log *zap.SugaredLogger, opts ...runtime.ServeMuxOption) {
	jsonpbMarshaler := &runtime.JSONPb{}
	jsonpbMarshaler.MarshalOptions = protojson.MarshalOptions{EmitUnpopulated: true}

	extraMutexOpts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption("application/json", jsonpbMarshaler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpbMarshaler),
	}

	if extraMutexOpts != nil {
		opts = append(opts, extraMutexOpts...)
	}

	gatewayMux := runtime.NewServeMux(opts...)
	gatewayServer := &http.Server{
		Addr:         httpPort,
		Handler:      allowCORS(gatewayMux),
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	dialAddr := fmt.Sprintf("127.0.0.1%s", grpcPort)
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ctx := context.Background()
			if err := api.RegisterTestServicesHandlerFromEndpoint(ctx, gatewayMux, dialAddr, dialOpts); err != nil {
				return err
			}

			log.Info("Starting gateway server.", zap.String("http_port", httpPort))
			go gatewayServer.ListenAndServe()
			log.Info("============ Gateway server started. ============ ")

			return nil
		},
		OnStop: func(ctx context.Context) error {

			err := gatewayServer.Shutdown(ctx)
			if err == nil {
				log.Info("============ Gateway server stopped.============ ")
			}
			return err
		},
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
	return
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
