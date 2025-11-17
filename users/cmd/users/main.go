package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"users/internal/users/app"
)

func main() {
	cfg := app.LoadConfig()

	server, err := app.NewHTTPServer(cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Println("user service started")

	// Graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	if err := server.Close(); err != nil {
		log.Printf("error shutting down server: %v", err)
	}

	log.Println("user service stopped")
}
