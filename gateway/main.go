package main

import (
	"context"
	"grpc-crud/pkg/pb/userpb"
	"log"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a cancelable context for the gateway lifecycle.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the HTTP request multiplexer for the gRPC Gateway.
	mux := runtime.NewServeMux()

	// Use an insecure connection to the local gRPC server for development.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Register the generated gateway handler from the service definitions.
	err := userpb.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:50051",
		opts,
	)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// Start the HTTP server on port 8080 and forward REST calls to gRPC.
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("REST Gateway running on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
