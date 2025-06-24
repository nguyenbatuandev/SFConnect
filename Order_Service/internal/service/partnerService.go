package service

import (
	"Order_Service/internal/entity"
	_interface "Order_Service/internal/interface"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type PartnerService struct {
	orderRepository _interface.OrderRepository
	authService     _interface.AuthService
	partnerRepository _interface.PartnerRepository
}

func NewPartnerService(orderRepository _interface.OrderRepository, authService _interface.AuthService, partnerRepository _interface.PartnerRepository) *PartnerService {
	return &PartnerService{
		orderRepository:  orderRepository,
		authService:      authService,
		partnerRepository: partnerRepository,
	}
}

func (s *PartnerService) GetAllOrdersByPartnerID(partnerID uuid.UUID, status string) ([]entity.OrderItemsRespone, error) {
	return s.partnerRepository.GetAllOrdersByPartnerID(partnerID, status)
}

func (s *PartnerService) UpdateOrderStatus(role entity.UserRole, orderItemID uuid.UUID) (*entity.OrderItemsRespone, error) {
	if role != entity.PartnerRole {
		return nil, errors.New("only partner can update order status")
	}

	status, err := s.orderRepository.GetOrderItemStatusByID(orderItemID)
	if err != nil {
		return nil, err
	}

	statusLower := strings.ToLower(strings.TrimSpace(status))
	switch statusLower {
	case "pending":
		order, err := s.orderRepository.UpdateOrderStatus(orderItemID, "confirmed")
		if err != nil {
			return nil, err
		}

		buyerID, err := s.orderRepository.GetBuyerIDByOrderItemID(orderItemID)
		if err == nil {
			go s.NotifyBuyerOrderReady(buyerID, order.OrderID)
		} else {
			fmt.Printf("Failed to get buyerID for notification: %v\n", err)
		}

		return order, nil

	default:
		return nil, fmt.Errorf("cannot update order with status: %s", status)
	}
}

func (s *PartnerService) NotifyBuyerOrderReady(buyerID uuid.UUID, orderID uuid.UUID) {
	notification := entity.Notification{
		BuyerID: buyerID,
		Title:   "Đơn hàng của bạn đã sẵn sàng",
		Message: "Đơn hàng của bạn đã sẵn sàng để nhận. Vui lòng xác nhận.",
		OrderID: orderID,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Lỗi khi mã hóa thông báo thành JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8083/api/notify", bytes.NewReader(body))
	if err != nil {
		fmt.Println("Lỗi khi tạo yêu cầu thông báo HTTP:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Lỗi khi gửi thông báo:", err)
		return
	}
	defer resp.Body.Close()
}