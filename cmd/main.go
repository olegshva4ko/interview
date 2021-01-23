package main

import (
	"context"
	"interview/configs"
	"interview/pkg/server"
	"log"
	"os"
	"os/signal"
)

func main() {
	config := configs.MakeConfig()

	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		osCall := <-interruption
		log.Printf("System call: %+v", osCall)
		cancel()
	}()

	if err := server.StartServer(ctx, config); err != nil {
		panic(err)
	}
}

