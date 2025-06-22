package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "partner"
	RoleGuest UserRole = "buyer"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      UserRole  `json:"role" gorm:"not null;default:'buyer'" validate:"required,oneof=admin partner buyer"`
	IsActive  bool      `json:"is_active" gorm:"not null;default:false"`
	CreatAt   time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}