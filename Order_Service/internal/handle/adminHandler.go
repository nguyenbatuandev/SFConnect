package handle

import (
	"Order_Service/internal/entity"
	_interface "Order_Service/internal/interface"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	orderService _interface.OrderServiceAdmin
	CommissionService _interface.PartnerCommissionService
}

func NewAdminHandler(orderService _interface.OrderServiceAdmin, 	CommissionService _interface.PartnerCommissionService) *AdminHandler {
	return &AdminHandler{
		orderService: orderService,
		CommissionService: CommissionService,
	}
}

func (a *AdminHandler) GetAllOrders(c *gin.Context) {
	orders, err := a.orderService.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

func (a *AdminHandler) GetOrderByID(c *gin.Context) {
	orderIDRaw := c.Param("orderId")

	orderID, err := uuid.Parse(orderIDRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid order ID"})
		return
	}

	order, err := a.orderService.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to retrieve order"})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, entity.ErrorResponse{Error: "Order not found"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Order retrieved successfully",
		Data:    order,
	})
}

func (a *AdminHandler) GetCommissionByOrderItemID(c *gin.Context) {
	orderItemIDRaw := c.Param("orderItemId")

	orderItemID, err := uuid.Parse(orderItemIDRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid order item ID"})
		return
	}

	commission, err := a.CommissionService.GetCommissionByOrderItemID(orderItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to retrieve commission"})
		return
	}

	if commission == nil {
		c.JSON(http.StatusNotFound, entity.ErrorResponse{Error: "Commission not found"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Commission retrieved successfully",
		Data:    commission,
	})
}

func (a *AdminHandler) GetCommissionByPartnerID(c *gin.Context) {
	orderItemIDRaw := c.Param("patrnerId")

	orderItemID, err := uuid.Parse(orderItemIDRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid order item ID"})
		return
	}

	commissions, err := a.CommissionService.GetCommissionByPartnerID(orderItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to retrieve commissions"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Commissions retrieved successfully",
		Data:    commissions,
	})
}
 