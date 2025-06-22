package handle

import (
	"Order_Service/internal/entity"
	_interface "Order_Service/internal/interface"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BuyerHandler struct {
	authService  _interface.AuthService
	orderService _interface.OrderServiceBuyer
}

func NewBuyerHandler(authService _interface.AuthService, orderService _interface.OrderServiceBuyer) *BuyerHandler {
	return &BuyerHandler{
		authService:  authService,
		orderService: orderService,
	}
}

func (b *BuyerHandler) GetAllOrdersByBuyerID(c *gin.Context) {
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
	status := c.Query("status") // sẽ trả về chuỗi rỗng nếu không có

	orders, err := b.orderService.GetAllOrdersByBuyerID(id, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

func (b *BuyerHandler) CreateOrder(c *gin.Context) {
	var orderRequest entity.CreateOrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid input", Message: err.Error()})
		return
	}

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

	// Tạo đối tượng Order
	order := entity.Order{
		ID:         uuid.New(),
		BuyerID:    id,
		Address:    orderRequest.Address,
		OrderItems: make([]entity.OrderItems, 0, len(orderRequest.OrderItems)),
	}

	// Chuyển đổi từ request items sang entity items
	for _, itemRequest := range orderRequest.OrderItems {
		productID, err := uuid.Parse(itemRequest.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid product ID", Message: err.Error()})
			return
		}

		partnerID, err := uuid.Parse(itemRequest.PartnerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid partner ID", Message: err.Error()})
			return
		}

		// Kiểm tra số lượng sản phẩm hợp lý
		if itemRequest.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Số lượng sản phẩm phải lớn hơn 0"})
			return
		}

		if itemRequest.Quantity > 100000 {
			c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Số lượng sản phẩm quá lớn", Message: fmt.Sprintf("Sản phẩm ID %s có số lượng %d vượt quá giới hạn cho phép (tối đa 100.000)", itemRequest.ProductID, itemRequest.Quantity)})
			return
		}

		order.OrderItems = append(order.OrderItems, entity.OrderItems{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ProductID: productID,
			PartnerID: partnerID,
			Quantity:  itemRequest.Quantity,
			Status:    "pending",
		})
	}

	// In ra ID để kiểm tra trong log khi lỗi xảy ra
	fmt.Printf("HANDLER: Đang cố gắng tạo Order với ID: %s\n", order.ID.String())

	newOrder, err := b.orderService.CreateOrder(&order)
	if err != nil {
		// Log lỗi chi tiết hơn nếu có thể
		fmt.Printf("HANDLER: Lỗi khi tạo đơn hàng: %v\n", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entity.SuccessResponse{
		Message: "Order created successfully",
		Data:    newOrder,
	})
}
func (b *BuyerHandler) UpdateOrderStatus(c *gin.Context) {
	// Lấy user_id từ context
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.ErrorResponse{Error: "User not authenticated"})
		return
	}
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Lấy role từ context
	roleRaw, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.ErrorResponse{Error: "User role not found"})
		return
	}
	role, ok := roleRaw.(entity.UserRole)
	if !ok || role != entity.BuyerRole {
		c.JSON(http.StatusForbidden, entity.ErrorResponse{Error: "You are not authorized to update order status"})
		return
	}

	// Parse order_id từ URL
	orderID, err := uuid.Parse(c.Param("orderItemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid order ID"})
		return
	}

	// Gọi service để cập nhật trạng thái
	order, err := b.orderService.UpdateOrderStatus(role, userID, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Error:   "Failed to update order status",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Order status updated successfully",
		Data:    order,
	})
}
