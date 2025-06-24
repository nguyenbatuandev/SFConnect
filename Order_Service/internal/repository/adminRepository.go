package repository

import (
	"Order_Service/internal/entity"
	"Order_Service/internal/interface"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db          *gorm.DB
	CallService _interface.CallService
}

func NewAdminRepository(db *gorm.DB, callService _interface.CallService) *AdminRepository {
	return &AdminRepository{
		db:          db,
		CallService: callService,
	}
}
func (r *AdminRepository) GetAllOrders() ([]entity.Order, error) {
		var orders []entity.Order
	if err := r.db.Preload("OrderItems").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *AdminRepository) GetOrderByID(orderID uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.Preload("OrderItems").First(&order, "id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}