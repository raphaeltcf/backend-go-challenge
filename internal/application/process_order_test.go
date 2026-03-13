package application

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/raphaeltcf/backend-go-challenge/internal/application/dto"
	"github.com/raphaeltcf/backend-go-challenge/internal/domain"
	"github.com/raphaeltcf/backend-go-challenge/internal/infra/metrics"
)

type mockRepository struct{}

func (m mockRepository) Save(ctx context.Context, order domain.Order) error {
	return nil
}

func (m mockRepository) FindByID(ctx context.Context, id string) (domain.Order, error) {
	return domain.Order{}, nil
}

func TestProcessOrder_ValidOrder(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	repo := mockRepository{}
	m := metrics.NewMetrics()
	useCase := NewProcessOrderUseCase(repo, logger, m)

	input := dto.OrderInputDTO{
		OrderID:   "order-123",
		UserID:    "user-456",
		Amount:    100.0,
		Timestamp: time.Now(),
	}

	output, err := useCase.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.OrderID != input.OrderID {
		t.Errorf("expected order ID %s, got %s", input.OrderID, output.OrderID)
	}
	if output.Status != string(domain.StatusProcessed) {
		t.Errorf("expected status %s, got %s", domain.StatusProcessed, output.Status)
	}

}

func TestProcessOrder_InvalidOrder(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	repo := mockRepository{}
	m := metrics.NewMetrics()
	useCase := NewProcessOrderUseCase(repo, logger, m)

	input := dto.OrderInputDTO{
		OrderID:   "",
		UserID:    "user-abc",
		Amount:    0,
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error for invalid order, got nil")
	}
}

func TestProcessOrder_NegativeAmount(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	repo := mockRepository{}
	m := metrics.NewMetrics()
	useCase := NewProcessOrderUseCase(repo, logger, m)

	input := dto.OrderInputDTO{
		OrderID:   "order-789",
		UserID:    "user-def",
		Amount:    -50.0,
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error for negative amount, got nil")
	}
}

func TestProcessOrder_EmptyUserID(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	repo := mockRepository{}
	m := metrics.NewMetrics()
	useCase := NewProcessOrderUseCase(repo, logger, m)
	input := dto.OrderInputDTO{
		OrderID:   "order-456",
		UserID:    "",
		Amount:    25.0,
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error for empty user ID, got nil")
	}
}

func TestProcessOrder_ZeroAmount(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	repo := mockRepository{}
	m := metrics.NewMetrics()
	useCase := NewProcessOrderUseCase(repo, logger, m)

	input := dto.OrderInputDTO{
		OrderID:   "order-321",
		UserID:    "user-ghi",
		Amount:    0.0,
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error for zero amount, got nil")
	}
}
