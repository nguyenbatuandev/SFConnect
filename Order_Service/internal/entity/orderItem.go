package entity

import (
	"github.com/google/uuid"
)

type OrderItems struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	PartnerID uuid.UUID `json:"partner_id" gorm:"type:uuid;not null"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`                   
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null;default:0.00"` 
	Status    string    `json:"status" gorm:"type:varchar(20);default:'pending'"`      
}
