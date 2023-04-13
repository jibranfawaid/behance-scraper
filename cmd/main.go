package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"optimus/config"
	"optimus/pkg"
	application "optimus/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Load config from .env
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	err = pkg.NewLogger()
	if err != nil {
		log.Fatalf("Failed to load logger config: %v", err)
	}

	tp, err := pkg.NewTracer()
	if err != nil {
		log.Error("Failed to load tracer config: " + err.Error())
	}

	pw, err := pkg.NewPlaywright()
	if err != nil {
		log.Error("Failed to start playwright: " + err.Error())
	}

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Error("Jaeger shutdown error: ", err)
		}
	}()

	application.RunServer(ctx, pw)
}
