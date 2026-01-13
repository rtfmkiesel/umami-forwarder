package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rtfmkiesel/umami-forwarder/umami"
)

func main() {
	config, err := umami.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded config '%#v'\n", config)

	umamiClient, err := umami.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	httpSrv := &http.Server{
		Addr:              ":8080",
		Handler:           umamiClient.Forward(),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	go func() {
		log.Printf("Starting forwarder server on :8080\n")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error while running forwarder server: %s", err)
		}
	}()

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)
	<-chanSignal

	log.Printf("Shutting down forwarder server...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatalf("Error while stopping forwarder server: %s\n", err)
	}
}
