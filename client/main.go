package main

import (
	"context"
	"flag"
	"go-client/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

var (
	addr = flag.String("addr", "localhost:8001", "the address to connect ")
	name = flag.String("name", "world", "Name to greet")
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(client()))
	})
	http.ListenAndServe(":8000", nil)
}

func init() {

}

func client() string {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
