package service

import (
	"Order_Service/internal/entity"
	_interface "Order_Service/internal/interface"
	"github.com/google/uuid"
)

type AdminService struct {
	orderRepository _interface.OrderRepository
	CommissionService _interface.PartnerCommissionService
	adminRepository _interface.AdminRepository
}

func NewAdminService(orderRepository _interface.OrderRepository,CommissionService _interface.PartnerCommissionService, adminRepository _interface.AdminRepository) *AdminService {
	return &AdminService{
		orderRepository:  orderRepository,
		CommissionService: CommissionService,
		adminRepository:  adminRepository,
	}
}

func (s *AdminService) GetAllOrders() ([]entity.Order, error) {
	orders, err := s.adminRepository.GetAllOrders()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *AdminService) GetOrderByID(orderID uuid.UUID) (*entity.Order, error) {
	order, err := s.adminRepository.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *AdminService) GetCommissionByOrderItemID(orderItemID uuid.UUID) (*entity.PartnerCommissionRespone, error) {
	commission, err := s.CommissionService.GetCommissionByOrderItemID(orderItemID)
	if err != nil {
		return nil, err
	}
	return commission, nil
}