package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/marciomarquesdesouza/go-rate-limiter/configs"
	"github.com/marciomarquesdesouza/go-rate-limiter/internal/infra/database/redis"
	"github.com/marciomarquesdesouza/go-rate-limiter/internal/infra/web"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	//Here you can change to another database if you like
	repository := redis.NewLimiterInfoRepository("localhost:6379", "", 0)

	server := web.NewServer(config.MaxRequestsPerSecond, config.BlockingTimeSeconds, repository)
	router := server.CreateServer()

	go func() {
		log.Println("Starting server on port", "8080")
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	// Create a timeout context for the graceful shutdown
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
