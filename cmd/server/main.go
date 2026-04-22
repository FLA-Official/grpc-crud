package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	profilepb "grpc-crud/gen/profile/v1"
	userpb "grpc-crud/gen/user/v1"
	"grpc-crud/internal/config"
	middleware "grpc-crud/internal/middlewares"
	profileHandler "grpc-crud/internal/profile/handler"
	profileRepo "grpc-crud/internal/profile/repo"
	profileService "grpc-crud/internal/profile/service"
	"grpc-crud/internal/user/handler"
	userRepo "grpc-crud/internal/user/repo"
	"grpc-crud/internal/user/service"
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
	userRepo := userRepo.NewUserRepo(db)
	profileRepo := profileRepo.NewProfileRepo(db)
	userService := service.NewUserService(userRepo, profileRepo)
	userHandler := handler.NewUserHandler(userService)

	// Build profile service layers.
	profileSvc := profileService.NewProfileService(profileRepo)
	profileHdlr := profileHandler.NewProfileHandler(profileSvc)

	// Create gRPC server and register the generated service implementations.
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			middleware.AuthInterceptor(cfg.JWTSecretKey),
		),
	)
	userpb.RegisterUserServiceServer(server, userHandler)
	profilepb.RegisterProfileServiceServer(server, profileHdlr)

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
