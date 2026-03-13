package queue

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/raphaeltcf/backend-go-challenge/internal/application/dto"
)

type MemoryQueue struct {
	ch chan dto.OrderInputDTO
}

func NewMemoryQueue(size int) *MemoryQueue {
	return &MemoryQueue{
		ch: make(chan dto.OrderInputDTO, size),
	}
}

func (q *MemoryQueue) Enqueue(ctx context.Context, order dto.OrderInputDTO) error {
	select {
	case q.ch <- order:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("failed to enqueue order: %w", ctx.Err())
	}
}

func (q *MemoryQueue) Dequeue(ctx context.Context) (dto.OrderInputDTO, error) {
	select {
	case order := <-q.ch:
		return order, nil
	case <-ctx.Done():
		return dto.OrderInputDTO{}, fmt.Errorf("failed to dequeue order: %w", ctx.Err())
	}
}

func StartGenerator(ctx context.Context, queue *MemoryQueue, logger *slog.Logger) {
	go func() {
		counter := 0
		for {
			select {
			case <-ctx.Done():
				logger.Info("Order generator stopped")
				return
			default:
				counter++

				var order dto.OrderInputDTO

				if counter%5 == 0 {
					order = dto.OrderInputDTO{
						OrderID:   "",
						UserID:    fmt.Sprintf("user-%d", time.Now().UnixNano()),
						Amount:    -1,
						Timestamp: time.Now(),
					}
					logger.Info("Generating invalid order")
				} else {
					order = dto.OrderInputDTO{
						OrderID:   fmt.Sprintf("order-%d", time.Now().UnixNano()),
						UserID:    fmt.Sprintf("user-%d", time.Now().UnixNano()),
						Amount:    float64(time.Now().UnixNano()%1000 + 1),
						Timestamp: time.Now(),
					}
				}

				if err := queue.Enqueue(ctx, order); err != nil {
					logger.Error("Failed to enqueue order", "error", err)
					return
				}
				logger.Info("Order enqueued", "order_id", order.OrderID)

				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}
