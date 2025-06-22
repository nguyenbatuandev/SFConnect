package _interface

import (
	"github.com/google/uuid"
	"user_service/internal/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByID(id uuid.UUID) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uuid.UUID) error
	ToggleUserLockByAdmin(id uuid.UUID, isActive bool) error
	GetAllUsersByAdmin() ([]*entity.User, error)
}

type AuthSerice interface {
	GenerateToken(user *entity.User) (string, error)
	ValidateToken(token string) (*Claims, error)
	HashPassword(password string) (string ,error)
	CheckPasswordHash(password, hash string) (bool)
}

type UserService interface {
	RegisterUser(user *entity.User) error
	LoginUser(user *entity.User) (*entity.User, error)
	GetProfile(id uuid.UUID) (*entity.User, error)
	UpdateProfile(id uuid.UUID, name string) (*entity.User, error)
	DeleteUser(id uuid.UUID) error
	ToggleUserLockByAdmin(id uuid.UUID) error
	GetAllUsersByAdmin() ([]*entity.User, error)
}

type Claims struct {
	UserID uuid.UUID       `json:"user_id"`
	Email  string          `json:"email"`
	Role   entity.UserRole `json:"role"`
}
