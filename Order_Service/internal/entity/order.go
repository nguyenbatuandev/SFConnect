package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	BuyerID      uuid.UUID    `json:"buyer_id" gorm:"type:uuid;not null"`
	CreateTime   time.Time    `json:"create_time" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	TotalPrice   float64      `json:"total_price" gorm:"type:decimal(10,2);default:0.00"`
	Address      string       `json:"address" gorm:"type:varchar(255);not null"`
	OrderItems   []OrderItems `json:"order_items" gorm:"constraint:OnDelete:CASCADE"`
}
