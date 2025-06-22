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

type UpdateRequest struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price      float64 `json:"price,omitempty" `
	Image      string  `json:"image,omitempty"`
}

type GetProductByNameResponse struct {
	Id		 uuid.UUID `json:"id"`
	Name        string  `json:"name" example:"Product Name"`
	Description string  `json:"description" example:"Product Description"`
	Price      float64 `json:"price" example:"99.99"`
	Image      string  `json:"image" example:"http://example.com/image.jpg"`
}

type UserRole string
type Claims struct {
	UserID uuid.UUID       `json:"user_id"`
	Email  string          `json:"email"`
	Role   UserRole       `json:"role"`
}

const (
	AdminRole    UserRole = "admin"
	BuyerRole    UserRole = "buyer"
	PartnerRole  UserRole = "partner"
)