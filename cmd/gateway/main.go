package main

import (
	"context"
	"log"
	"net/http"
	"time"

	profilepb "grpc-crud/gen/profile/v1"
	userpb "grpc-crud/gen/user/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a cancelable context for the gateway lifecycle.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the HTTP request multiplexer for the gRPC Gateway.
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		if key == "Authorization" {
			return "authorization", true
		}
		return key, false
	}),
	)

	// Use an insecure connection to the local gRPC server for development.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Register the generated gateway handler for UserService.
	err := userpb.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:50051",
		opts,
	)
	if err != nil {
		log.Fatalf("failed to register user service gateway: %v", err)
	}

	// Register the generated gateway handler for ProfileService.
	err = profilepb.RegisterProfileServiceHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:50051",
		opts,
	)
	if err != nil {
		log.Fatalf("failed to register profile service gateway: %v", err)
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
