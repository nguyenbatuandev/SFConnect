package entity

import (
	"time"

	"github.com/google/uuid"
)

type PartnerCommission struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	PartnerID uuid.UUID `json:"partner_id" gorm:"type:uuid;not null"`
	OrderItems uuid.UUID `json:"order_items" gorm:"type:uuid;not null"`
	CommissionRate float64 `json:"commission_rate" gorm:"type:decimal(10,2);not null"`
	CommissionAmount float64 `json:"commission_amount" gorm:"type:decimal(10,2);not null"`
	CreateTime time.Time `json:"create_time" gorm:"type:timestamp"`
}