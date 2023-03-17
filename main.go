package main

import (
	"log"
	"net"

	"github.com/yu-yk/median-svc/median"
	"github.com/yu-yk/median-svc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	medianServer := median.NewServer()

	pb.RegisterMedianServer(grpcServer, medianServer)
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
