package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/canyouhearthemusic/todo-list/internal/routes"
	log "github.com/sirupsen/logrus"
)

func main() {
	var wg sync.WaitGroup

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: routes.New(),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Printf("Starting Application on :%v port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :%v: %v\n", port, err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-signalCh

	log.Warnln("\nGracefully shutting down HTTP server...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Errorf("HTTP server shutdown error: %s", err)
	} else {
		log.Infoln("HTTP server shut down gracefully.")
	}

	wg.Wait()

	log.Infoln("Shutdown complete.")
}
