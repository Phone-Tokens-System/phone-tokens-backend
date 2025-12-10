package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	server "phone-tokens/internal/adapter/in"
	"phone-tokens/internal/app"
)

// Entry point for the monolithic HTTP server.
func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	userSvc, tokenSvc, err := app.BuildService(cfg)
	if err != nil {
		log.Fatalf("failed to initialize sms_service: %v", err)
	}

	httpServer, err := server.NewHTTPServer(cfg.HTTPPort, cfg.JWTSecret, userSvc, tokenSvc)
	if err != nil {
		log.Fatalf("failed to initialize HTTP server: %v", err)
	}

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
