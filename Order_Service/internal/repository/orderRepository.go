package repository

import (
	"Order_Service/internal/entity"
	_interface "Order_Service/internal/interface"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
	CallService _interface.CallService
}

func NewOrderRepository(db *gorm.DB, callService _interface.CallService) *postgresRepository {
	return &postgresRepository{db: db, CallService: callService}
}

func (r *postgresRepository) GetAllOrderByID(orderID uuid.UUID, status string) (*entity.Order, error) {
	var order entity.Order
	query := r.db.Preload("OrderItems").Where("id = ?", orderID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
} //


func (r *postgresRepository) UpdateOrderStatus(orderItemID uuid.UUID, status string) (*entity.OrderItemsRespone, error) {
	var orderItem entity.OrderItems

	// Cập nhật trạng thái của order item cụ thể
	if err := r.db.Model(&entity.OrderItems{}).
		Where("id = ?", orderItemID).
		Update("status", status).Error; err != nil {
		return nil, err
	}

	// Lấy lại order item mới nhất sau khi cập nhật
	if err := r.db.First(&orderItem, "id = ?", orderItemID).Error; err != nil {
		return nil, err
	}

	product, err := r.CallService.GetProductByID(orderItem.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product info for ProductID %s: %w", orderItem.ProductID, err)
	}

	orderItenRp := entity.OrderItemsRespone{
		ID:        orderItem.ID,
		OrderID:   orderItem.OrderID,
		Product:   *product,
		PartnerID: orderItem.PartnerID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		Status:    orderItem.Status,
	}

	return &orderItenRp, nil
}

func (r *postgresRepository) CheckOwnerOrder(buyerID, orderItemID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.
		Model(&entity.OrderItems{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("order_items.id = ? AND orders.buyer_id = ?", orderItemID, buyerID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, errors.New("you are not the owner of this order item")
	}
	return true, nil
}

func (r *postgresRepository) GetOrderItemStatusByID(orderItemID uuid.UUID) (string, error) {
	var status string
	err := r.db.Model(&entity.OrderItems{}).
		Select("status").
		Where("id = ?", orderItemID).
		Take(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}

func (r *postgresRepository) GetOrderItemByID(orderItemID uuid.UUID) (*entity.OrderItemsRespone, error) {
	var orderItem entity.OrderItems
	if err := r.db.First(&orderItem, "id = ?", orderItemID).Error; err != nil {
		return nil, err
	}

	// Lấy thông tin sản phẩm
	product, err := r.CallService.GetProductByID(orderItem.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product info for ProductID %s: %w", orderItem.ProductID, err)
	}

	// Mapping sang response
	opRp := &entity.OrderItemsRespone{
		ID:        orderItem.ID,
		OrderID:   orderItem.OrderID,
		Product:   *product,
		PartnerID: orderItem.PartnerID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		Status:    orderItem.Status,
	}

	return opRp, nil
}

func (r *postgresRepository) GetBuyerIDByOrderItemID(orderItemID uuid.UUID) (uuid.UUID, error) {
	var orderItem entity.OrderItems
	if err := r.db.Select("order_id").First(&orderItem, "id = ?", orderItemID).Error; err != nil {
		return uuid.Nil, err
	}

	var order entity.Order
	if err := r.db.Select("buyer_id").First(&order, "id = ?", orderItem.OrderID).Error; err != nil {
		return uuid.Nil, err
	}

	return order.BuyerID, nil
}