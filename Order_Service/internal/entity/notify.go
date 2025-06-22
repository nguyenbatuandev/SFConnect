package entity

import "github.com/google/uuid"

type Notification struct {
	BuyerID uuid.UUID   `json:"buyer_id"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	OrderID uuid.UUID `json:"order_id"`
}