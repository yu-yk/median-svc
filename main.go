package main

import (
	"context"
	"log"
	"net"

	"github.com/yu-yk/median-svc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedMedianServer
}

func (s *server) PushNumber(ctx context.Context, req *pb.PushNumberRequest) (*pb.PushNumberResponse, error) {
	return &pb.PushNumberResponse{
		Status: &pb.Status{},
	}, nil
}

func (s *server) GetMedianRequest(ctx context.Context, req *pb.GetMedianRequest) (*pb.GetMedianResponse, error) {
	return &pb.GetMedianResponse{}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	medianServer := &server{}

	pb.RegisterMedianServer(grpcServer, medianServer)
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
