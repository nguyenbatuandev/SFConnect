package entity

import (
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid input"`
	Message string `json:"message,omitempty" example:"Validation failed"`
}

type SuccessResponse struct {
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

type UserRole string
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   UserRole  `json:"role"`
}

const (
	AdminRole   UserRole = "admin"
	BuyerRole   UserRole = "buyer"
	PartnerRole UserRole = "partner"
)

type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	PartnerID string `json:"partner_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type CreateOrderRequest struct {
	Address    string             `json:"address" binding:"required"`
	OrderItems []OrderItemRequest `json:"order_items" binding:"required"`
}
