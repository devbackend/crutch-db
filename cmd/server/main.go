package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devbackend/crutch-db/internal/operation"
	"github.com/devbackend/crutch-db/internal/parser"
	"github.com/devbackend/crutch-db/internal/server"
	"github.com/devbackend/crutch-db/internal/storage"
)

func main() {
	port := flag.Int("port", 2357, "running port")

	flag.Parse()

	srv := server.New(
		server.WithPort(*port),
		server.WithParser(
			parser.New(),
		),
		server.WithOperationRunner(
			operation.NewRunner(
				operation.WithStorage(
					storage.New(),
				),
			),
		),
	)

	go func(srv server.Server) {
		err := srv.Start()
		if err != nil {
			log.Fatal(err)
		}
	}(srv)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	log.Println("Stopping...")

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	err := srv.Stop(timeout)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Stopped")
}
