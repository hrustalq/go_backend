package main

import (
	"log"
	"net"

	"github.com/hrustalq/go_backend/internal"
	"github.com/hrustalq/go_backend/proto/auth"
	"google.golang.org/grpc"
)

func main() {
	db, err := internal.ConnectDatabase("user=postgres password=password dbname=authdb sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	redisClient, err := internal.ConnectRedis("localhost:6379", "", 0)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	authService := &internal.AuthService{
		DB:        db,
		Redis:     redisClient,
		JWTSecret: []byte("your-secret-key"),
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, authService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
