package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/raphaeltcf/backend-go-challenge/internal/application"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/logger"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/metrics"
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

	m := metrics.NewMetrics()
	useCase := application.NewProcessOrderUseCase(repo, log, m)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := queue.NewMemoryQueue(100)

	queue.StartGenerator(ctx, q, log)

	wp := worker.NewWorkerPool(q, useCase, log)
	wp.Start(ctx, 5)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"total_processed": m.TotalProcessedOrders,
			"total_failed":    m.TotalFailedOrders,
			"total_invalid":   m.TotalInvalidOrders,
			"failure_rate":    m.FailureRate(),
		})
	})

	go http.ListenAndServe(":8080", nil)

	log.Info("Order processing system started")
	log.Info("metrics available at http://localhost:8080/metrics")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	<-sigchan
	log.Info("Shutting down order processing system")
	cancel()
}
