package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"grpc-crud/internal/config"
	"grpc-crud/internal/user/handler"
	"grpc-crud/internal/user/repo"
	"grpc-crud/internal/user/service"
	"grpc-crud/pkg/pb/userpb"
)

func main() {
	// Load database configuration from environment variables.
	cfg := config.GetConfig()

	// Build PostgreSQL DSN using the loaded configuration.
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Connect to PostgreSQL using sqlx.
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Build the application layers: repository, service, and gRPC handler.
	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Create gRPC server and register the generated service implementation.
	server := grpc.NewServer()
	userpb.RegisterUserServiceServer(server, userHandler)

	// Listen for incoming gRPC connections on port 50051.
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server running at :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
