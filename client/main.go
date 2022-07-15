package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-k8s-example/pb"
	"log"
	"net/http"
	"time"
)

var (
	addr       = flag.String("addr", "localhost:8001", "the address to connect ")
	name       = flag.String("name", "world", "Name to greet")
	clientPort = flag.String("clientPort", ":8000", "client listen port")
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(client()))
	})
	http.ListenAndServe(*clientPort, nil)
}

func client() string {
	flag.Parse()
	address := fmt.Sprintf("dns:///%s", *addr)
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("invoke error: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
	return r.GetMessage()
}
