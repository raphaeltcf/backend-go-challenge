package queue

import (
	"context"
	"fmt"

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
