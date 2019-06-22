package main

import (
	"context"
	"flag"
	"log"
	"net"

	pb "server/api/v1"

	"google.golang.org/grpc"
)

// Server represents the gRPC server
type Server struct {
}

// SayHello generates response to a Ping request
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "你好！世界！"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &Server{})
	if err := grpcServer.Serve(lis);err !=  nil{
		panic(err)
	}
}