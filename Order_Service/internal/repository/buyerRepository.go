package repository

import (
	"Order_Service/internal/entity"
	"Order_Service/internal/interface"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BuyerRepository struct{
	db *gorm.DB
	CallService _interface.CallService
}

func NewBuyerRepository(db *gorm.DB, callService _interface.CallService) *BuyerRepository {
	return &BuyerRepository{
		db:          db,
		CallService: callService,
	}
}

func (r *BuyerRepository) CreateOrder(order *entity.Order) (*entity.Order, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Đảm bảo rollback khi có lỗi hoặc panic
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
		// Rollback nếu có bất kỳ lỗi nào trong transaction
		if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// Đảm bảo ID Order được tạo mới nếu chưa có.
	// Sử dụng clause.Insert{} để đảm bảo GORM thực hiện INSERT thuần túy.
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	if err := tx.Clauses(clause.Insert{}).Create(order).Error; err != nil {
		tx.Error = fmt.Errorf("failed to create order: %w", err)
		return nil, tx.Error
	}

	var totalPrice float64 = 0
	orderItemsToPersist := make([]entity.OrderItems, 0, len(order.OrderItems))

	for i := range order.OrderItems {
		item := &order.OrderItems[i]

		// Khởi tạo ID mới cho Order Item và liên kết với Order chính
		item.ID = uuid.New()
		item.OrderID = order.ID
		item.Status = "pending" // Thiết lập trạng thái mặc định

		if item.Quantity <= 0 {
			tx.Error = errors.New("order item quantity must be greater than zero")
			return nil, tx.Error
		}

		// Lấy thông tin sản phẩm từ Product Service
		product, err := r.CallService.GetProductByID(item.ProductID)
		if err != nil {
			tx.Error = fmt.Errorf("failed to get product info for ProductID %s: %w", item.ProductID, err)
			return nil, tx.Error
		}
		if product == nil {
			tx.Error = fmt.Errorf("product with ID %s not found", item.ProductID)
			return nil, tx.Error
		}

		// Cập nhật giá sản phẩm cho OrderItem (đã giả định product.Price là uint64)
		item.Price = product.Price

		// Tính toán tổng giá cho mục hàng này.
		itemTotal := item.Price * float64(item.Quantity) // <-- Đây là phép tính đúng

		// Cộng vào tổng giá trị đơn hàng
		totalPrice += itemTotal

		// Thêm vào slice để insert hàng loạt
		orderItemsToPersist = append(orderItemsToPersist, *item)
	}

	// Thực hiện insert hàng loạt tất cả order items
	if len(orderItemsToPersist) > 0 {
		if err := tx.Create(&orderItemsToPersist).Error; err != nil {
			tx.Error = fmt.Errorf("failed to create order items: %w", err)
			return nil, tx.Error
		}
	}

	// Cập nhật tổng giá của đơn hàng chính
	order.TotalPrice = totalPrice
	if err := tx.Model(order).Update("total_price", totalPrice).Error; err != nil {
		tx.Error = fmt.Errorf("failed to update order total price: %w", err)
		return nil, tx.Error
	}

	// Commit transaction nếu tất cả thành công
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

func (r *BuyerRepository) GetAllOrdersByBuyerID(buyerID uuid.UUID, status string) ([]entity.OrderRespone, error) {
	var orders []entity.Order

	// preload OrderItems
	query := r.db.Preload("OrderItems").Where("buyer_id = ?", buyerID)
	if status != "" {
		// lấy order có item có status khớp
		subquery := r.db.Model(&entity.OrderItems{}).
			Select("order_id").
			Where("status = ?", status)
		query = query.Where("id IN (?)", subquery)
	}

	// lấy orders từ DB
	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	// map sang []OrderRespone (có product)
	var orderRespList []entity.OrderRespone

	for _, order := range orders {
		var orderItemsResp []entity.OrderItemsRespone

		for _, item := range order.OrderItems {
			//status filter
			if status != "" && item.Status != status {
				continue
			}

			// Lấy product info
			product, err := r.CallService.GetProductByID(item.ProductID)
			if err != nil {
				return nil, fmt.Errorf("failed to get product info for ProductID %s: %w", item.ProductID, err)
			}
			
			orderItemsResp = append(orderItemsResp, entity.OrderItemsRespone{
				ID:        item.ID,
				OrderID:   item.OrderID,
				Product:   *product,
				PartnerID: item.PartnerID,
				Quantity:  item.Quantity,
				Price:     item.Price,
				Status:    item.Status,
			})
		}

		orderRespList = append(orderRespList, entity.OrderRespone{
			ID:         order.ID,
			BuyerID:    order.BuyerID,
			CreateTime: order.CreateTime,
			TotalPrice: order.TotalPrice,
			Address:    order.Address,
			OrderItems: orderItemsResp,
		})
	}

	return orderRespList, nil
}