package repository

import (
	"Order_Service/internal/entity"
	_interface "Order_Service/internal/interface" // Import đúng interface package
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PartnerCommissionRepository struct {
	db              *gorm.DB
	orderRepository _interface.OrderRepository
}

// Constructor: inject OrderRepository vào
func NewPartnerCommissionRepository(db *gorm.DB, orderRepo _interface.OrderRepository) *PartnerCommissionRepository {
	return &PartnerCommissionRepository{
		db:              db,
		orderRepository: orderRepo,
	}
}

// Lấy hoa hồng theo partner
func (r *PartnerCommissionRepository) GetCommissionByPartnerID(partnerID uuid.UUID) (*[]entity.PartnerCommissionRespone, error) {
	var commissions []entity.PartnerCommission
	err := r.db.Where("partner_id = ?", partnerID).Find(&commissions).Error
	if err != nil {
		return nil, err
	}

	var responses []entity.PartnerCommissionRespone
	for _, commission := range commissions {
		orderItem, err := r.orderRepository.GetOrderItemByID(commission.OrderItems)
		if err != nil {
			return nil, err
		}

		responses = append(responses, entity.PartnerCommissionRespone{
			ID:               commission.ID,
			OrderItems:       *orderItem,
			CommissionRate:   commission.CommissionRate,
			CommissionAmount: commission.CommissionAmount,
			CreateTime:       commission.CreateTime,
		})
	}

	return &responses, nil
}

// Tạo bản ghi hoa hồng mới
func (r *PartnerCommissionRepository) CreateCommission(commission *entity.PartnerCommission) (*entity.PartnerCommission, error) {
	if err := r.db.Create(commission).Error; err != nil {
		return nil, err
	}
	return commission, nil
}

func (r *PartnerCommissionRepository) GetCommissionByOrderItemID(orderItemID uuid.UUID) (*entity.PartnerCommissionRespone, error) {
	var commission entity.PartnerCommission
	err := r.db.Where("order_items = ?", orderItemID).First(&commission).Error
	if err != nil {
		return nil, err
	}

	orderItem, err := r.orderRepository.GetOrderItemByID(commission.OrderItems)
	if err != nil {
		return nil, err
	}

	response := &entity.PartnerCommissionRespone{
		ID:               commission.ID,
		OrderItems:       *orderItem,
		CommissionRate:   commission.CommissionRate,
		CommissionAmount: commission.CommissionAmount,
		CreateTime:       commission.CreateTime,
	}

	return response, nil
}