package main

import (
	"log"
	"net"

	pb "github.com/ahmedalialphasquad123/calculationService/proto"
	"github.com/ahmedalialphasquad123/calculationService/server"
	"google.golang.org/grpc"
)

func main() {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 16), // 16MB max receive message size
		grpc.MaxSendMsgSize(1024 * 1024 * 16), // 16MB max send message size
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCalculationServiceServer(grpcServer, &server.Server{})

	log.Println("gRPC server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
