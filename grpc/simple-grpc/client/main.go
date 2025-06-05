package main

import (
	"context"
	"log"
	"time"

	pb "example.com/simple-grpc/hello"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("‼️ could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Herman"})
	if err != nil {
		log.Fatalf("‼️ error calling SayHello: %v", err)
	}

	log.Printf("✅ Got server response: %s", resp.Message)
}
