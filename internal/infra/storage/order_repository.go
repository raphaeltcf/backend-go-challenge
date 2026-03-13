package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/raphaeltcf/backend-go-challenge/internal/domain"
)

type SQLiteOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) SQLiteOrderRepository {
	return SQLiteOrderRepository{db: db}
}

func (r SQLiteOrderRepository) Save(ctx context.Context, order domain.Order) error {
	query := `
    INSERT INTO orders (id, user_id, amount, status, error, created_at, processed_at) 
    VALUES (?, ?, ?, ?, ?, ?, ?)
    ON CONFLICT(id) DO UPDATE SET 
        status=excluded.status,
        error=excluded.error,
        processed_at=excluded.processed_at
`
	_, err := r.db.ExecContext(ctx, query, order.ID, order.UserID, order.Amount, string(order.Status), "", order.CreatedAt, order.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r SQLiteOrderRepository) FindByID(ctx context.Context, id string) (domain.Order, error) {
	query := `SELECT id, user_id, amount, status, created_at FROM orders WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var order domain.Order
	var status string
	err := row.Scan(&order.ID, &order.UserID, &order.Amount, &status, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Order{}, fmt.Errorf("order not found: %w", err)
		}
		return domain.Order{}, fmt.Errorf("failed to get order by ID: %w", err)
	}

	order.Status = domain.OrderStatus(status)
	return order, nil
}
