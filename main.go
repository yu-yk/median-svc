package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yu-yk/median-svc/median"
	"github.com/yu-yk/median-svc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// gprc
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	medianServer := median.NewServer()

	proto.RegisterMedianServer(grpcServer, medianServer)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Panicf("failed to serve: %v", err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := proto.RegisterMedianHandlerFromEndpoint(ctx, mux, ":3000", opts); err != nil {
		log.Panicf("failed to serve: %v", err)
	}

	// http
	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
