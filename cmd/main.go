package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"phone-tokens/internal/adapter/in"
	httpadapter "phone-tokens/internal/adapter/in/http"
	"phone-tokens/internal/app"
	"syscall"
)

// Entry point for the monolithic HTTP server.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by your JWT token

func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	services, err := app.BuildService(cfg)
	if err != nil {
		log.Fatalf("failed to initialize sms: %v", err)
	}
	handlers := httpadapter.BuildHandlers(*services)

	httpServer, err := in.NewHTTPServer(cfg.HTTPPort, cfg.JWTSecret, *handlers)
	if err != nil {
		log.Fatalf("failed to initialize HTTP server: %v", err)
	}
	fmt.Println(services.Cert.Storage)
	go func() {
		log.Printf("HTTP server started on :%s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	log.Println("monolith started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	if err := httpServer.Close(); err != nil {
		log.Printf("error shutting down HTTP server: %v", err)
	}

	log.Println("monolith stopped")
}
