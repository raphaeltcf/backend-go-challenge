package worker

import (
	"context"
	"log/slog"

	"github.com/raphaeltcf/backend-go-challenge/internal/application"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/correlation"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/queue"
)

type WorkerPool struct {
	queue   *queue.MemoryQueue
	useCase application.ProcessOrderUseCase
	logger  *slog.Logger
}

func NewWorkerPool(q *queue.MemoryQueue, useCase application.ProcessOrderUseCase, logger *slog.Logger) WorkerPool {
	return WorkerPool{
		queue:   q,
		useCase: useCase,
		logger:  logger,
	}
}

func (wp WorkerPool) Start(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			wp.logger.Info("worker started", "worker_id", workerID)
			for {
				order, err := wp.queue.Dequeue(ctx)
				if err != nil {
					wp.logger.Info("worker stopped", "worker_id", workerID)
					return
				}

				orderCtx := correlation.WithCorrelationID(ctx, order.CorrelationID)

				_, err = wp.useCase.Execute(orderCtx, order)
				if err != nil {
					wp.logger.ErrorContext(orderCtx, "failed to process order",
						"worker_id", workerID,
						"order_id", order.OrderID,
						"correlation_id", order.CorrelationID,
						"error", err,
					)
				}
			}
		}(i)
	}
}
