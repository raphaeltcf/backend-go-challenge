package infra_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/raphaeltcf/backend-go-challenge/internal/application"
	"github.com/raphaeltcf/backend-go-challenge/internal/application/dto"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/metrics"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/queue"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/storage"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/worker"
)

func TestE2E_ProcessOrder(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := storage.NewConnection(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := storage.NewOrderRepository(db)
	m := metrics.NewMetrics()
	useCase := application.NewProcessOrderUseCase(repo, logger, m)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := queue.NewMemoryQueue(10)

	wp := worker.NewWorkerPool(q, useCase, logger)
	wp.Start(ctx, 2)

	input := dto.OrderInputDTO{
		OrderID:   "order-123",
		UserID:    "user-456",
		Amount:    100.0,
		Timestamp: time.Now(),
	}

	if err := q.Enqueue(ctx, input); err != nil {
		t.Fatalf("Failed to enqueue order: %v", err)
	}

	time.Sleep(1 * time.Second)

	order, err := repo.FindByID(ctx, input.OrderID)
	if err != nil {
		t.Fatalf("Failed to find order: %v", err)
	}

	if order.ID != input.OrderID {
		t.Errorf("Expected order ID %s, got %s", input.OrderID, order.ID)
	}
	if order.Status != "processed" {
		t.Errorf("Expected status 'processed', got '%s'", order.Status)
	}
}
