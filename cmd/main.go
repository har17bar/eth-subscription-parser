package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"txparser/internal/api/rest"
	"txparser/internal/client"
	"txparser/internal/job"
	"txparser/internal/service"
	"txparser/internal/storage"
)

const (
	apiKey = ""
)

func main() {
	inMemStore := storage.NewStorage()
	etCli := client.MustETHClient(apiKey)
	ethService := service.NewEth(inMemStore, etCli)
	userService := service.NewUser(inMemStore, ethService)
	subscriptionHandler := rest.NewSubscription(userService)
	ethHandler := rest.NewTransactions(ethService)
	server := rest.NewServer(subscriptionHandler, ethHandler)
	jobEth := job.NewEth(ethService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background job
	if err := jobEth.Start(ctx); err != nil {
		log.Fatalf("failed to start job: %v", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Handle OS signals for shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("Shutting down gracefully...")

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Println("Shutdown complete.")
}
