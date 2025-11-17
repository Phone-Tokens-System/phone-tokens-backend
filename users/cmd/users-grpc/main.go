package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	grpcadapter "users/internal/users/adapter/in/grpc"
	"users/internal/users/app"
	"users/internal/users/service/users"
)

func main() {
	cfg := app.LoadConfig()

	repo, err := app.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}
	svc := users.NewService(repo, users.Config{
		JWTSecret:       cfg.JWTSecret,
		JWTExpiresInSec: cfg.JWTExpiresInSec,
	})

	if cfg.GRPCPort == "" {
		log.Fatalf("GRPC_PORT is required")
	}

	addr := ":" + cfg.GRPCPort
	server, lis, err := grpcadapter.Start(addr, svc)
	if err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

	go func() {
		log.Printf("gRPC server started on %s", addr)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	server.GracefulStop()
	log.Println("gRPC server stopped")
}
