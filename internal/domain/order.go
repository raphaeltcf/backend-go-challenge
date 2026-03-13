package domain

import (
	"errors"
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusProcessed OrderStatus = "processed"
	StatusFailed    OrderStatus = "failed"
)

type Order struct {
	ID        string
	UserID    string
	Amount    float64
	Status    OrderStatus
	CreatedAt time.Time
}

func NewOrder(id, userID string, amount float64, createdAt time.Time) Order {
	return Order{
		ID:        id,
		UserID:    userID,
		Amount:    amount,
		Status:    StatusPending,
		CreatedAt: createdAt,
	}
}

func (o Order) Validate() error {
	if o.UserID == "" {
		return errors.New("order user id cannot be empty")
	}
	if o.Amount <= 0 {
		return errors.New("order amount must be greater than zero")
	}
	if o.ID == "" {
		return errors.New("order id cannot be empty")
	}
	return nil
}
