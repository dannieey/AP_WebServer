package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"AP_WebServer/internal/server"
	"AP_WebServer/internal/store"
	"AP_WebServer/internal/worker"
)

func main() {
	store := store.NewStore[string, string]()
	srv := server.NewServer(store)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: srv.Routes(),
	}

	stopWorker := make(chan struct{})
	go worker.StartWorker(stopWorker, srv)

	go func() {
		log.Println("Server started on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down...")
	close(stopWorker)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)

	log.Println("Server stopped")
}
