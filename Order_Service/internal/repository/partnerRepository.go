package repository

import (
	"Order_Service/internal/entity"
	"Order_Service/internal/interface"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PartnerRepository struct {
	db          *gorm.DB
	CallService _interface.CallService
}

func NewPartnerRepository(db *gorm.DB, callService _interface.CallService) *PartnerRepository {
	return &PartnerRepository{
		db:          db,
		CallService: callService,
	}
}

func (r *PartnerRepository) GetAllOrdersByPartnerID(partnerID uuid.UUID, status string) ([]entity.OrderItemsRespone, error) {
	var items []entity.OrderItems

	// Query cơ bản
	query := r.db.Where("partner_id = ?", partnerID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Thực thi truy vấn
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	// Tạo danh sách về
	var result []entity.OrderItemsRespone

	for _, item := range items {
		// Lấy thông tin sản phẩm
		product, err := r.CallService.GetProductByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product info for ProductID %s: %w", item.ProductID, err)
		}
		if product == nil {
			return nil, fmt.Errorf("product with ID %s not found", item.ProductID)
		}

		//Map
		result = append(result, entity.OrderItemsRespone{
			ID:        item.ID,
			OrderID:   item.OrderID,
			Product:   *product,
			PartnerID: item.PartnerID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Status:    item.Status,
		})
	}

	return result, nil
}