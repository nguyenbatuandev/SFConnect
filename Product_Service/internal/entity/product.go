package entity

import "github.com/google/uuid"

type Product struct {
	Id          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Description string    `json:"description" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	Image       string    `json:"image" gorm:"not null"`
}
