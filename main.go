package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yu-yk/median-svc/median"
	"github.com/yu-yk/median-svc/proto"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func newGRPCService(logger *zap.Logger) (*grpc.Server, *median.Server) {
	grpcServer := grpc.NewServer()
	medianServer := median.NewServer(logger)

	return grpcServer, medianServer
}

func registerGRPCService(lc fx.Lifecycle, grpcServer *grpc.Server, medianServer *median.Server) {
	proto.RegisterMedianServer(grpcServer, medianServer)
	reflection.Register(grpcServer)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", ":3000")
			if err != nil {
				return err
			}
			log.Println("serving grpc at tcp:3000")
			go grpcServer.Serve(listener)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})
}

func registerHTTPGateway(lc fx.Lifecycle, mediaServer *median.Server) *http.Server {
	mux := runtime.NewServeMux()
	server := &http.Server{Addr: ":3001", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := proto.RegisterMedianHandlerServer(ctx, mux, mediaServer)
			// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
			// err := proto.RegisterMedianHandlerFromEndpoint(ctx, mux, ":3000", opts)
			if err != nil {
				return err
			}
			log.Println("serving http at :3001")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Close()
		},
	})

	return server
}

func main() {
	fx.New(
		fx.Provide(
			newGRPCService,
			zap.NewExample,
		),
		fx.Invoke(registerGRPCService),
		fx.Invoke(registerHTTPGateway),
	).Run()
}
