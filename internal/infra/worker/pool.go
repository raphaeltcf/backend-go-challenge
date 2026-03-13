package worker

import (
	"context"
	"log/slog"

	"github.com/raphaeltcf/backend-go-challenge/internal/application"
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
				_, err = wp.useCase.Execute(ctx, order)
				if err != nil {
					wp.logger.Error("failed to process order", "worker_id", workerID, "error", err)
				}
			}
		}(i)
	}
}
