package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"grpc/api_pb/pb"
	service2 "grpc/service"
	"log"
	"net"
)

var (
	grpcServer = flag.String("grpc-server",  "127.0.0.1:8080", "gRPC server")
)

func main(){

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := grpc.NewServer()
	service := &service2.InnService{}

	pb.RegisterInnServiceServer(server, service)

	listener, err := net.Listen("tcp", *grpcServer)
	if err != nil{
		log.Fatalf("Fail to create gRPC listener: %s", err.Error())
	}
	if err = server.Serve(listener); err != nil{
		log.Fatalf("Fail to start server: %s", err.Error())
	}


}
