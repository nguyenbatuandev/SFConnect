package service

import (
	"Product_Service/internal/entity"
	_interface "Product_Service/internal/interface"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo  _interface.ProductRepository
	authService  _interface.AuthService
	cacheService _interface.CacheService
}

func NewProductService(productRepo _interface.ProductRepository, authService _interface.AuthService, cacheService _interface.CacheService) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		authService:  authService,
		cacheService: cacheService,
	}
}

func (s *ProductService) CreateProduct(product *entity.Product) (*entity.Product, error) {
	if err := s.productRepo.CreateProduct(product); err != nil {
		return nil, errors.New("failed to create product: " + err.Error())
	}
	cacheKey := fmt.Sprintf("product:%s", product.Id.String())
	if err := s.cacheService.Set(cacheKey, product, 60*time.Minute); err != nil {
		log.Printf("Failed to cache product %s: %v", product.Id.String(), err)
	} else {
		log.Printf("Successfully cached product with key: %s", cacheKey)
	}
	_ = s.cacheService.Delete("products:all")

	return product, nil
}

func (s *ProductService) GetProductByID(id uuid.UUID) (*entity.Product, error) {
	cacheKey := fmt.Sprintf("product:%s", id.String())
	var cachedProduct *entity.Product
	if s.cacheService.Exists(cacheKey) {
		if err := s.cacheService.Get(cacheKey, &cachedProduct); err == nil {
			return cachedProduct, nil
		}
	}

	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, errors.New("product not found: " + err.Error())
	}

	if err := s.cacheService.Set(cacheKey, product, 60*time.Minute); err != nil {
		log.Printf("Failed to cache product %s: %v", id.String(), err)
	}

	return product, nil
}

func (s *ProductService) UpdateProduct(id uuid.UUID, product *entity.UpdateRequest) (*entity.Product, error) {
	existingProduct, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, errors.New("product not found: " + err.Error())
	}

	oldName := existingProduct.Name

	if product.Name != "" {
		existingProduct.Name = product.Name
	}
	if product.Description != "" {
		existingProduct.Description = product.Description
	}
	if product.Price != 0 {
		existingProduct.Price = product.Price
	}
	if product.Image != "" {
		existingProduct.Image = product.Image
	}

	if err := s.productRepo.UpdateProduct(existingProduct); err != nil {
		return nil, errors.New("failed to update product: " + err.Error())
	}

	// Cập nhật cache
	productKey := fmt.Sprintf("product:%s", id.String())
	_ = s.cacheService.Set(productKey, existingProduct, 60*time.Minute)

	// Xóa cache theo tên cũ và danh sách sản phẩm
	nameKey := fmt.Sprintf("product_by_name:%s", oldName)
	_ = s.cacheService.Delete(nameKey)
	_ = s.cacheService.DeletePattern("products:all")

	return existingProduct, nil
}

func (s *ProductService) DeleteProduct(id uuid.UUID) error {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return errors.New("product not found: " + err.Error())
	}

	if err := s.productRepo.DeleteProduct(id); err != nil {
		return errors.New("failed to delete product: " + err.Error())
	}

	// Xóa cache cụ thể
	productKey := fmt.Sprintf("product:%s", id.String())
	nameKey := fmt.Sprintf("product_by_name:%s", product.Name)

	_ = s.cacheService.Delete(productKey)
	_ = s.cacheService.Delete(nameKey)
	_ = s.cacheService.DeletePattern("products:all")

	return nil
}

func (s *ProductService) GetProductByName(name string) ([]*entity.Product, error) {

	cacheKey := fmt.Sprintf("product_by_name:%s", name)
	var cachedProducts []*entity.Product

	if s.cacheService.Exists(cacheKey) {
		if err := s.cacheService.Get(cacheKey, &cachedProducts); err == nil {
			return cachedProducts, nil
		}
	}

	product, err := s.productRepo.GetProductByName(name)
	if err != nil {
		return nil, errors.New("product not found: " + err.Error())
	}

	if err := s.cacheService.Set(cacheKey, product, 60*time.Minute); err != nil {
		log.Printf("Failed to cache products: %v", err)
	}

	return product, nil
}

func (s *ProductService) GetAllProducts() ([]*entity.Product, error) {
	cacheKey := "products:all"
	var products []*entity.Product

	if s.cacheService.Exists(cacheKey) {
		if err := s.cacheService.Get(cacheKey, &products); err == nil {
			return products, nil
		}
	}

	products, err := s.productRepo.GetAllProducts()
	if err != nil {
		return nil, errors.New("failed to retrieve products: " + err.Error())
	}

	_ = s.cacheService.Set(cacheKey, products, 60*time.Minute)

	return products, nil
}
