package handle

import (
	"Order_Service/internal/entity"
	_interface "Order_Service/internal/interface"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PartnerHandler struct {
	authService       _interface.AuthService
	orderService      _interface.OrderServicePartner
	partnerCommission _interface.PartnerCommissionService
}

func NewPartnerHandler(authService _interface.AuthService, orderService _interface.OrderServicePartner, partnerCommission _interface.PartnerCommissionService) *PartnerHandler {
	return &PartnerHandler{
		authService:       authService,
		orderService:      orderService,
		partnerCommission: partnerCommission,
	}
}

func (b *PartnerHandler) GetAllOrdersByPartnerID(c *gin.Context) {
	buyerIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.ErrorResponse{Error: "User not authenticated"})
		return
	}

	id, ok := buyerIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Lấy status từ query string
	status := c.Param("status") // sẽ trả về chuỗi rỗng nếu không có

	orders, err := b.orderService.GetAllOrdersByPartnerID(id, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

func (b *PartnerHandler) UpdateOrderStatus(c *gin.Context) {
	orderIDRaw := c.Param("orderItemId")

	orderID, err := uuid.Parse(orderIDRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid order item ID"})
		return
	}

	roleRaw, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.ErrorResponse{Error: "User role not provided"})
		return
	}

	role, ok := roleRaw.(entity.UserRole)
	if !ok {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Invalid user role"})
		return
	}

	order, err := b.orderService.UpdateOrderStatus(role, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Order status updated successfully",
		Data:    order,
	})
}

func (b *PartnerHandler) GetCommissionByPartnerID(c *gin.Context) {
	partnerIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.ErrorResponse{Error: "User not authenticated"})
		return
	}

	id, ok := partnerIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	commissions, err := b.partnerCommission.GetCommissionByPartnerID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to retrieve commissions"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Commissions retrieved successfully",
		Data:    commissions,
	})
}
