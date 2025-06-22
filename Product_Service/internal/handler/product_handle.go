package handler

import (
	"Product_Service/internal/entity"
	_interface "Product_Service/internal/interface"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService _interface.ProductService
	authService    _interface.AuthService
}

func NewProductHandler(productService _interface.ProductService, authService _interface.AuthService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		authService:    authService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
		return
	}
	product.Id = uuid.New()
	createdProduct, err := h.productService.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entity.SuccessResponse{
		Message: "Product created successfully",
		Data:    createdProduct,
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	// Lấy ID sản phẩm từ tham số URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	// Gọi service để lấy sản phẩm theo ID
	product, err := h.productService.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Nếu không tìm thấy sản phẩm, trả về lỗi NotFound
	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")


	// Parse product ID từ path
	productID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	// Bind dữ liệu cập nhật
	var updateReq entity.UpdateRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
		return
	}

	// Gọi service để cập nhật
	updatedProduct, err := h.productService.UpdateProduct(productID, &updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: err.Error()})
		return
	}

	// Trả về kết quả thành công
	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Gọi service để lấy danh sách sản phẩm
	products, err := h.productService.GetAllProducts()
	if err != nil {
		// Nếu có lỗi trong quá trình lấy sản phẩm, trả về lỗi 500 với thông báo lỗi
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: err.Error()})
		return
	}

	// Tạo slice để chứa response
	productsResponse := make([]entity.GetProductByNameResponse, len(products))
	for i, p := range products {
		productsResponse[i] = entity.GetProductByNameResponse{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Image:       p.Image,
		}
	}

	// Trả về kết quả thành công kèm dữ liệu sản phẩm
	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Products retrieved successfully",
		Data:    productsResponse,
	})
}

func (h *ProductHandler) GetProductByName(c *gin.Context) {
	// Lấy tên sản phẩm từ tham số URL
	name := c.Param("name")
	// Nếu tên sản phẩm rỗng, trả về lỗi BadRequest
	if name == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Product name is required"})
		return
	}

	// Gọi service tìm sản phẩm theo tên
	product, err := h.productService.GetProductByName(name)
	if err != nil {
		// Nếu lỗi trong quá trình tìm kiếm, trả về lỗi InternalServerError
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: err.Error()})
		return
	}

	// Nếu không tìm thấy sản phẩm nào, trả về lỗi NotFound
	if product == nil {
		c.JSON(http.StatusNotFound, entity.ErrorResponse{Error: "Product not found"})
		return
	}

	// Chuyển đổi dữ liệu sản phẩm sang định dạng response
	productResponse := make([]entity.GetProductByNameResponse, len(product))
	for i, p := range product {
		productResponse[i] = entity.GetProductByNameResponse{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Image:       p.Image,
		}
	}

	// Trả về kết quả thành công với danh sách sản phẩm tìm được
	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Product found",
		Data:    productResponse,
	})
}
