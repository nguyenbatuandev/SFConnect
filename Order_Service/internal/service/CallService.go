package service

import (
	"Order_Service/internal/entity"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CallService struct {
}

func NewCallService() *CallService {
	return &CallService{}
}

func (s *CallService) GetListPartner(authHeader string) ([]entity.GetPartnerResponse, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "http://user_service:8080/api/buyer/get-list-partner", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user-service error: %s", string(body))
	}

	var wrapper entity.GetPartnerResponseWrapper
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return wrapper.Data, nil
}

func (s *CallService) PushJWT(c *fiber.Ctx) error {
	 // Lấy token từ request gốc
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization"})
    }

    // Gọi sang service khác
    partners, err := s.GetListPartner(authHeader)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(partners)
}

func (r *CallService) GetProductByID(productID uuid.UUID) (*entity.Product, error) {
	// Xây dựng URL
	url := fmt.Sprintf("http://product_service:8082/api/v1/product-id/%s", productID.String())

	// Tạo HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo HTTP request cho sản phẩm: %w", err)
	}

	// Gửi request bằng http.DefaultClient (hoặc tạo mới client tại đây)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi gửi request tới product service: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi đọc response body từ product service: %w", err)
	}

	// Kiểm tra mã trạng thái
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product service trả về trạng thái lỗi: %d - %s", resp.StatusCode, string(body))
	}

	// Parse JSON
	var productResp entity.ProductResponse
	if err := json.Unmarshal(body, &productResp); err != nil {
		return nil, fmt.Errorf("lỗi khi parse response sản phẩm: %w", err)
	}

	// Kiểm tra sản phẩm hợp lệ
	if productResp.Data.Id == uuid.Nil {
		return nil, fmt.Errorf("không tìm thấy sản phẩm với ID %s trong phản hồi từ service", productID)
	}

	return &productResp.Data, nil
}

