package application

import (
	"context"

	"github.com/raphaeltcf/backend-go-challenge/internal/application/dto"
	"github.com/raphaeltcf/backend-go-challenge/internal/domain"
)

type OrderRepository interface {
	Save(ctx context.Context, order domain.Order) error
	FindByID(ctx context.Context, id string) (domain.Order, error)
}

type OrderQueue interface {
	Enqueue(ctx context.Context, order dto.OrderInputDTO) error
	Dequeue(ctx context.Context) (dto.OrderInputDTO, error)
}
