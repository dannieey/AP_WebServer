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
	// Создаём generic Store[string,string]
	store := store.NewStore[string, string]()
	srv := server.NewServer(store)

	// чекаем изменение ключей +2
	srv.Store().Set("k1", "v1")
	srv.Store().Set("k2", "v2")

	// Запуск фонового worker-а
	stopWorker := make(chan struct{})
	go worker.StartWorker(srv, stopWorker)

	// HTTP сервер с handler от нашего сервера
	httpSrv := &http.Server{
		Addr:    ":8080",
		Handler: srv.Routes(), // srv.Routes() возвращает http.Handler
	}

	// Запуск сервера в отдельной горутине
	go func() {
		log.Println("Server started on :8080")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Ловим OS сигналы для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop // ждём сигнал

	log.Println("Shutting down server...")

	// Останавливаем worker
	close(stopWorker)

	// Graceful shutdown HTTP сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}

	log.Println("Server exited gracefully")

}
