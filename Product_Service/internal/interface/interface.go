package _interface

import (
	"Product_Service/internal/entity"
	"time"

	"github.com/google/uuid"
)

type ProductRepository interface {
	CreateProduct(product *entity.Product) error
	GetProductByID(id uuid.UUID) (*entity.Product, error)
	UpdateProduct(product *entity.Product) error
	DeleteProduct(id uuid.UUID) error
	GetProductByName(name string) ([]*entity.Product, error)
	GetAllProducts() ([]*entity.Product, error)
}

type AuthService interface {
	ValidateToken(token string) (*entity.Claims, error)
}

type ProductService interface {
	CreateProduct(product *entity.Product) (*entity.Product, error)
	GetProductByID(id uuid.UUID) (*entity.Product, error)
	UpdateProduct(id uuid.UUID, product *entity.UpdateRequest) (*entity.Product, error)
	DeleteProduct(id uuid.UUID) error
	GetProductByName(name string) ([]*entity.Product, error)
	GetAllProducts() ([]*entity.Product, error)
}

type CacheService interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Exists(key string) bool
	DeletePattern(pattern string) error
}
