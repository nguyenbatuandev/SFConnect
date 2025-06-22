package service

import (
	"Order_Service/internal/entity"
	_interface "Order_Service/internal/interface"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BuyerService struct {
	authService              _interface.AuthService
	orderRepository          _interface.OrderRepository
	PartnerCommissionService _interface.PartnerCommissionService
}

func NewBuyerService(authService _interface.AuthService, orderRepository _interface.OrderRepository, PartPartnerCommissionService _interface.PartnerCommissionService) *BuyerService {
	return &BuyerService{
		authService:              authService,
		orderRepository:          orderRepository,
		PartnerCommissionService: PartPartnerCommissionService,
	}
}

func (s *BuyerService) CreateOrder(order *entity.Order) (*entity.Order, error) {
	return s.orderRepository.CreateOrder(order)
}

func (s *BuyerService) GetAllOrdersByBuyerID(buyerID uuid.UUID, status string) ([]entity.OrderRespone, error) {
	return s.orderRepository.GetAllOrdersByBuyerID(buyerID, status)
}

func (s *BuyerService) UpdateOrderStatus(role entity.UserRole, buyerID, orderItemID uuid.UUID) (*entity.OrderItemsRespone, error) {
	// check role
	if role != entity.BuyerRole {
		return nil, errors.New("only buyer can update order status")
	}

	// check quyền sở hữu order
	isOwner, err := s.orderRepository.CheckOwnerOrder(buyerID, orderItemID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("you are not the owner of this order")
	}

	// Lấy status hiện tại
	currentStatus, err := s.orderRepository.GetOrderItemStatusByID(orderItemID)
	if err != nil {
		return nil, err
	}

	statusLower := strings.ToLower(strings.TrimSpace(currentStatus))

	switch statusLower {
	case "pending":
		// Cho phép hủy đơn
		return s.orderRepository.UpdateOrderStatus(orderItemID, "cancel")

	case "confirmed":
		order, err := s.orderRepository.GetOrderItemByID(orderItemID)
		if err != nil {
			return nil, err
		}
		if order == nil {
			return nil, errors.New("order not found")
		}
		commissionRate := 0.1 // 10%

		partnerCommission := &entity.PartnerCommission{
			ID:               uuid.New(),
			PartnerID:        order.PartnerID,
			OrderItems:       orderItemID,
			CommissionRate:   commissionRate,
			CommissionAmount: order.Price * commissionRate,
			CreateTime:       time.Now(),
		}

		createdCommission, err := s.PartnerCommissionService.CreateCommission(partnerCommission)
		if err != nil {
			return nil, fmt.Errorf("failed to create partner commission: %w", err)
		}
		if createdCommission == nil {
			return nil, errors.New("failed to create partner commission")
		}

		return s.orderRepository.UpdateOrderStatus(orderItemID, "complete")

	default:
		return nil, fmt.Errorf("cannot update order with status: %s", currentStatus)
	}
}
