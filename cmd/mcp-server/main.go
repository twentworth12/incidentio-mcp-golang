package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tomwentworth/incidentio-mcp-golang/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	srv := server.New()
	srv.Start(ctx)
}