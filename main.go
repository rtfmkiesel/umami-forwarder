package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/rtfmkiesel/kisslog"

	"github.com/rtfmkiesel/umami-forwarder/umami"
)

var version string = "@DEV" // Adjusted by the Makefile
const addr string = ":8080" // Listening addr, hardcoded for now

func main() {
	if err := logger.InitDefault("umami-forwarder" + version); err != nil {
		panic(err)
	}

	log := logger.New("main")

	config, err := umami.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	umamiClient, err := umami.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           umamiClient.Forward(),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	go func() {
		log.Info("Starting forwarder server on %s\n", addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("error while running forwarder server: %s", err)
		}
	}()

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)
	<-chanSignal

	log.Info("Shutting down forwarder server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("error while stopping forwarder server: %s", err)
	}
}
