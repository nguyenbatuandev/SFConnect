package _interface

import (
	"Order_Service/internal/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthService interface {
	ValidateToken(token string) (*entity.Claims, error)
}

type OrderServiceBuyer interface {
	CreateOrder(order *entity.Order) (*entity.Order, error)
	GetAllOrdersByBuyerID(buyerID uuid.UUID, status string) ([]entity.OrderRespone, error)
	UpdateOrderStatus(role entity.UserRole, buyerID uuid.UUID, orderID uuid.UUID) (*entity.OrderItemsRespone, error)
}

type OrderServicePartner interface {
	GetAllOrdersByPartnerID(partnerID uuid.UUID, status string) ([]entity.OrderItemsRespone, error)
	UpdateOrderStatus(role entity.UserRole, orderID uuid.UUID) (*entity.OrderItemsRespone, error)
}

type OrderServiceAdmin interface {
	GetAllOrders() ([]entity.Order, error)
	GetOrderByID(orderID uuid.UUID) (*entity.Order, error)
	GetCommissionByOrderItemID(orderItemID uuid.UUID) (*entity.PartnerCommissionRespone, error)
}

type OrderRepository interface {
	CreateOrder(order *entity.Order) (*entity.Order, error)
	GetOrderByID(orderID uuid.UUID) (*entity.Order, error)
	GetAllOrdersByBuyerID(buyerID uuid.UUID, status string) ([]entity.OrderRespone, error)
	GetAllOrdersByPartnerID(partnerID uuid.UUID, status string) ([]entity.OrderItemsRespone, error)
	UpdateOrderStatus(orderID uuid.UUID, status string) (*entity.OrderItemsRespone, error)
	GetAllOrders() ([]entity.Order, error)
	CheckOwnerOrder(buyerID, orderID uuid.UUID) (bool, error)
	GetOrderItemStatusByID(orderItemID uuid.UUID) (string, error)
	GetOrderItemByID(orderItemID uuid.UUID) (*entity.OrderItemsRespone, error)
	GetAllOrderByID(orderID uuid.UUID, status string) (*entity.Order, error)
	GetBuyerIDByOrderItemID(orderItemID uuid.UUID) (uuid.UUID, error)
}

type PartnerCommissionRepository interface {
	GetCommissionByPartnerID(partnerID uuid.UUID) (*[]entity.PartnerCommissionRespone, error)
	CreateCommission(commission *entity.PartnerCommission) (*entity.PartnerCommission, error)
	GetCommissionByOrderItemID(orderItemID uuid.UUID) (*entity.PartnerCommissionRespone, error)
}

type PartnerCommissionService interface {
	GetCommissionByPartnerID(partnerID uuid.UUID) (*[]entity.PartnerCommissionRespone, error)
	CreateCommission(commission *entity.PartnerCommission) (*entity.PartnerCommission, error)
	GetCommissionByOrderItemID(orderItemID uuid.UUID) (*entity.PartnerCommissionRespone, error)

}

type ProductService interface {
	GetProductByID(productID uuid.UUID) (*entity.Product, error)
}

type CallService interface {
	GetListPartner(authHeader string) ([]entity.GetPartnerResponse, error)
	GetProductByID(productID uuid.UUID) (*entity.Product, error)
	PushJWT(c *fiber.Ctx) error
}