package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/raphaeltcf/backend-go-challenge/internal/application/dto"
	"github.com/raphaeltcf/backend-go-challenge/internal/domain"
)

type ProcessOrderUseCase struct {
	repo   OrderRepository
	logger *slog.Logger
}

func NewProcessOrderUseCase(repo OrderRepository, logger *slog.Logger) ProcessOrderUseCase {
	return ProcessOrderUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc ProcessOrderUseCase) Execute(ctx context.Context, input dto.OrderInputDTO) (dto.OrderOutputDTO, error) {
	uc.logger.InfoContext(ctx, "receiving order", "order_id", input.OrderID)

	order := domain.NewOrder(input.OrderID, input.UserID, input.Amount, input.Timestamp)

	if err := order.Validate(); err != nil {
		uc.logger.ErrorContext(ctx, "invalid order", "order_id", input.OrderID, "error", err)
		return dto.OrderOutputDTO{
			OrderID:     input.OrderID,
			Status:      string(domain.StatusFailed),
			Error:       err.Error(),
			ProcessedAt: time.Now(),
		}, fmt.Errorf("invalid order: %w", err)
	}
	var processErr error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			uc.logger.WarnContext(ctx, "retrying order processing", "order_id", input.OrderID, "attempt", attempt)
			time.Sleep(2 * time.Second)
		}

		processErr = simulateProcessing(order)
		if processErr == nil {
			break
		}
	}

	if processErr != nil {
		uc.logger.ErrorContext(ctx, "failed to process order after retries", "order_id", input.OrderID, "error", processErr)
		order.Status = domain.StatusFailed
		uc.repo.Save(ctx, order)
		return dto.OrderOutputDTO{
			OrderID:     input.OrderID,
			Status:      string(domain.StatusFailed),
			Error:       processErr.Error(),
			ProcessedAt: time.Now(),
		}, fmt.Errorf("failed to process order after retries: %w", processErr)
	}

	order.Status = domain.StatusProcessed
	if err := uc.repo.Save(ctx, order); err != nil {
		uc.logger.ErrorContext(ctx, "failed to save processed order", "order_id", input.OrderID, "error", err)
		return dto.OrderOutputDTO{}, fmt.Errorf("failed to save order: %w", err)
	}

	uc.logger.InfoContext(ctx, "order processed successfully", "order_id", input.OrderID)
	return dto.OrderOutputDTO{
		OrderID:     input.OrderID,
		Status:      string(domain.StatusProcessed),
		ProcessedAt: time.Now(),
	}, nil
}

func simulateProcessing(order domain.Order) error {
	return nil
}
