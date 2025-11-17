package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	grpcadapter "users/internal/users/adapter/in/grpc"
	"users/internal/users/app"
)

func main() {
	cfg := app.LoadConfig()

	// Общая инициализация доменного сервиса (БД, миграции, репозиторий).
	svc, err := app.BuildService(cfg)
	if err != nil {
		log.Fatalf("failed to initialize service: %v", err)
	}

	// HTTP‑сервер.
	httpServer, err := app.NewHTTPServer(cfg, svc)
	if err != nil {
		log.Fatalf("failed to initialize HTTP server: %v", err)
	}

	// gRPC‑сервер.
	if cfg.GRPCPort == "" {
		log.Fatal("GRPC_PORT is required")
	}
	grpcAddr := ":" + cfg.GRPCPort
	grpcServer, grpcLis, err := grpcadapter.Start(grpcAddr, svc)
	if err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

	// Запуск HTTP и gRPC параллельно.
	go func() {
		log.Printf("HTTP server started on :%s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		log.Printf("gRPC server started on %s", grpcAddr)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	log.Println("user service started (HTTP + gRPC)")

	// Graceful shutdown для обоих серверов.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	if err := httpServer.Close(); err != nil {
		log.Printf("error shutting down HTTP server: %v", err)
	}

	grpcServer.GracefulStop()

	log.Println("user service stopped")
}
