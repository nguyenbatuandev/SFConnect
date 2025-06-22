package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	Message string  `json:"message"`
	Data    Product `json:"data"`
}

type Product struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64    `json:"price"`
	Image       string    `json:"image"`
}

type OrderItemsRespone struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	Product   Product   `json:"product"`
	PartnerID uuid.UUID `json:"partner_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
}

type OrderRespone struct {
	ID           uuid.UUID    `json:"id"`
	BuyerID      uuid.UUID    `json:"buyer_id"`
	CreateTime   time.Time    `json:"create_time"`
	TotalPrice   float64      `json:"total_price"`
	Address      string       `json:"address"`
	OrderItems   []OrderItemsRespone `json:"order_items"`
}

type PartnerCommissionRespone struct {
	ID uuid.UUID `json:"id"`
	OrderItems OrderItemsRespone `json:"order_items"`
	CommissionRate float64 `json:"commission_rate"`
	CommissionAmount float64 `json:"commission_amount"`
	CreateTime time.Time `json:"create_time"`
}