package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/raphaeltcf/backend-go-challenge/internal/application/dto"
	"github.com/raphaeltcf/backend-go-challenge/internal/domain"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/correlation"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/metrics"
)

type ProcessOrderUseCase struct {
	repo    OrderRepository
	logger  *slog.Logger
	metrics *metrics.Metrics
}

func NewProcessOrderUseCase(repo OrderRepository, logger *slog.Logger, m *metrics.Metrics) ProcessOrderUseCase {
	return ProcessOrderUseCase{
		repo:    repo,
		logger:  logger,
		metrics: m,
	}
}

func (uc ProcessOrderUseCase) Execute(ctx context.Context, input dto.OrderInputDTO) (dto.OrderOutputDTO, error) {
	corrID := correlation.FromContext(ctx)

	uc.logger.InfoContext(ctx, "receiving order",
		uc.logFields(input.OrderID, "pending", corrID, nil)...)

	order := domain.NewOrder(input.OrderID, input.UserID, input.Amount, input.Timestamp)

	if err := order.Validate(); err != nil {
		uc.logger.ErrorContext(ctx, "invalid order",
			uc.logFields(input.OrderID, string(domain.StatusFailed), corrID, err)...)
		uc.metrics.IncInvalidOrders()
		return dto.OrderOutputDTO{
			OrderID:     input.OrderID,
			Status:      string(domain.StatusFailed),
			Error:       err.Error(),
			ProcessedAt: time.Now(),
		}, fmt.Errorf("invalid order: %w", err)
	}
	var processErr error
	existing, err := uc.repo.FindByID(ctx, input.OrderID)
	if err == nil && existing.Status == domain.StatusProcessed {
		uc.logger.InfoContext(ctx, "order already processed, skipping",
			uc.logFields(existing.ID, string(existing.Status), corrID, nil)...)
		return dto.OrderOutputDTO{
			OrderID:     existing.ID,
			Status:      string(existing.Status),
			ProcessedAt: time.Now(),
		}, nil
	}
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			uc.logger.WarnContext(ctx, "retrying order processing",
				append(uc.logFields(input.OrderID, "retrying", corrID, nil),
					"attempt", attempt)...)
			time.Sleep(2 * time.Second)
		}

		processErr = simulateProcessing(order)
		if processErr == nil {
			break
		}
	}

	if processErr != nil {
		uc.logger.ErrorContext(ctx, "failed to process order after retries",
			uc.logFields(input.OrderID, string(domain.StatusFailed), corrID, processErr)...)
		uc.metrics.IncFailedOrders()
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
		uc.logger.ErrorContext(ctx, "failed to save processed order",
			uc.logFields(input.OrderID, string(domain.StatusFailed), corrID, err)...)
		return dto.OrderOutputDTO{}, fmt.Errorf("failed to save order: %w", err)
	}

	uc.logger.InfoContext(ctx, "order processed successfully",
		uc.logFields(input.OrderID, string(domain.StatusProcessed), corrID, nil)...)
	uc.metrics.IncProcessedOrders()
	return dto.OrderOutputDTO{
		OrderID:     input.OrderID,
		Status:      string(domain.StatusProcessed),
		ProcessedAt: time.Now(),
	}, nil
}

func simulateProcessing(order domain.Order) error {
	return nil
}

func (uc ProcessOrderUseCase) logFields(orderID, status, corrID string, err error) []any {
	fields := []any{
		"order_id", orderID,
		"status", status,
		"correlation_id", corrID,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	}
	if err != nil {
		fields = append(fields, "error", err.Error())
	}
	return fields
}
