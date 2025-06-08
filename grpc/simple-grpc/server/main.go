package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "example.com/simple-grpc/hello" // pb - alias to access protobuff generated code
	"google.golang.org/grpc"           // imports grpc framework
)

// Represents the Server
type helloServer struct {
	// Go does not use classic Inheritance, uses composition and interfaces.
	// "extending" (embeding) protobuf to override "SayHello"
	// Embeding: "I want helloServer to have all `pb.UnimplementedHelloServiceServer` methods and fields"
	pb.UnimplementedHelloServiceServer
}

// How the server handles a "SayHello" call.
//   - `(s *helloServer)` - method receiver. SayHello is an associated method to `*helloServer` (Struct).
//     Where `*` indicates a pointer to helloServer and is named as `s` as convention (e.g. it's like self / this).
//     Java: `class HelloServer { String sayHello(Context ctx, ClientRequest clientRequest) { ... } }`
//   - `(*pb.HelloResponse, error)`- return values.
//     Where `*pb.HelloResponse` is the server response.
//     Java: `HelloResponse sayHello(...) throws Exception { ... }`
func (s *helloServer) SayHello(ctx context.Context, clientRequest *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("✅ Request Received: '%s'", clientRequest.Name)
	return &pb.HelloResponse{Message: fmt.Sprintf("Hello, %s!", clientRequest.Name)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("‼️ failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, &helloServer{})

	log.Println("✅ gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("‼️ failed to serve: %v", err)
	}
}
