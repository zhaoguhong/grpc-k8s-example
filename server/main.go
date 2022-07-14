package main

import (
	"context"
	"flag"
	"fmt"
	"go-client/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	port = flag.Int("port", 8001, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedHelloServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	hostname, _ := os.Hostname()
	// include hostname, use test loadBalance
	return &pb.HelloResponse{Message: "Hello " + in.GetName() + ", server hostname: " + hostname}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &server{}) // Register server
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
