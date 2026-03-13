package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/raphaeltcf/backend-go-challenge/internal/application"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/logger"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/queue"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/storage"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/worker"
)

func main() {
	log := logger.NewLogger()

	db, err := storage.NewConnection("orders.sqlite")
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := storage.NewOrderRepository(db)

	useCase := application.NewProcessOrderUseCase(repo, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := queue.NewMemoryQueue(100)

	queue.StartGenerator(ctx, q, log)

	wp := worker.NewWorkerPool(q, useCase, log)
	wp.Start(ctx, 5)

	log.Info("Order processing system started")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	<-sigchan
	log.Info("Shutting down order processing system")
	cancel()
}
