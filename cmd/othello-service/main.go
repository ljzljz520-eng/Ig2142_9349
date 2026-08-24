package main

import (
	"context"
	"example.com/othello-records/internal/app"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	databasePath := flag.String("db", "othello.db", "sqlite database path")
	address := flag.String("addr", ":8080", "http listen address")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	services, err := app.OpenServices(ctx, *databasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer services.Close()
	if err := services.Run(ctx, *address); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}
