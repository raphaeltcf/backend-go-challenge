package dto

import (
	"time"
)

type OrderOutputDTO struct {
	OrderID     string    `json:"order_id"`
	Status      string    `json:"status"`
	Error       string    `json:"error,omitempty"`
	ProcessedAt time.Time `json:"processed_at"`
}
