package queue

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/raphaeltcf/backend-go-challenge/internal/application/dto"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/correlation"
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

func generateCorrelationID() string {
	return fmt.Sprintf("corr-%d", time.Now().UnixNano())
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

				corrID := generateCorrelationID()
				orderCtx := correlation.WithCorrelationID(ctx, corrID)

				var order dto.OrderInputDTO

				if counter%5 == 0 {
					order = dto.OrderInputDTO{
						OrderID:       "",
						UserID:        fmt.Sprintf("user-%d", time.Now().UnixNano()),
						Amount:        -1,
						Timestamp:     time.Now(),
						CorrelationID: corrID,
					}
					logger.InfoContext(orderCtx, "generating invalid order", "correlation_id", corrID)
				} else {
					order = dto.OrderInputDTO{
						OrderID:       fmt.Sprintf("order-%d", time.Now().UnixNano()),
						UserID:        fmt.Sprintf("user-%d", time.Now().UnixNano()),
						Amount:        float64(time.Now().UnixNano()%1000 + 1),
						Timestamp:     time.Now(),
						CorrelationID: corrID,
					}
				}

				if err := queue.Enqueue(orderCtx, order); err != nil {
					logger.Error("Failed to enqueue order", "error", err)
					return
				}
				logger.InfoContext(orderCtx, "order enqueued",
					"order_id", order.OrderID,
					"correlation_id", corrID,
				)

				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}
