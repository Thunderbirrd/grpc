package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"grpc/api_pb/pb"
	"log"
	"net/http"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint",  "127.0.0.1:8080", "gRPC server endpoint address")
	port = flag.Int("port",  8081, "port for generated endpoints")
)

func main(){
	var conn *grpc.ClientConn
	conn, err :=grpc.Dial(":8080", grpc.WithInsecure())

	if err != nil{
		log.Fatalf("could not connect: %s", err.Error())
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %s", err.Error())
		}
	}(conn)

	c := pb.NewInnServiceClient(conn)

	var inns = []string{"545655303053", "5017096885"}
	for _, inn := range inns{
		response, err := c.GetInfoByInn(context.Background(), &pb.InnRequest{Inn: inn})
		if err != nil{
			log.Fatalf("Error: %s", err)
		}else if response != nil && response.GetName() != ""{
			log.Printf("Response from server: %s", response)
		}else{
			log.Printf("Empty response from server")
		}
	}

	if err := startHttpClient(); err != nil{
		log.Fatalf("Error whith HTTP client: %s", err.Error())
	}

}

func startHttpClient() error{
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	options := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterInnServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, options)
	if err != nil{
		log.Fatalf("Error while creating HTTP client: %s", err.Error())
		return err
	}

	url := fmt.Sprintf(":%d", *port)
	err = http.ListenAndServe(url, mux)
	if err != nil {
		log.Fatalf("Error while starting HTTP client: %s", err.Error())
		return err
	}
	return nil
}
