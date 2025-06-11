package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnelD/eventbus/bus"
	"github.com/AnelD/eventbus/filewatcher"
	"github.com/AnelD/eventbus/ws"
)

func main() {
	// Initialize the event bus
	eventBus := bus.New()

	// Start HTTP server with WebSocket support
	mux := ws.NewRouter(eventBus)

	// Start event bus logging
	go eventBus.LogAll()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Start file watcher
	shutDownWatcher := make(chan struct{})
	go filewatcher.PublishFileEventsWS("./watched_folder", shutDownWatcher)

	<-stop // Wait for a signal

	log.Println("Shutting down server...")
	close(shutDownWatcher)
	// Create a deadline to wait for current operations to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}
