package repository

import (
	"Product_Service/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetAllProducts() ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *postgresRepository) CreateProduct(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *postgresRepository) GetProductByID(id uuid.UUID) (*entity.Product, error) {
	returnedProduct := &entity.Product{}
	if err := r.db.First(returnedProduct, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return returnedProduct, nil
}

func (r *postgresRepository) UpdateProduct(product *entity.Product) error {
	return r.db.Save(product).Error
}

func (r *postgresRepository) DeleteProduct(id uuid.UUID) error {
	return r.db.Delete(&entity.Product{}, "id = ?", id).Error
}

func (r *postgresRepository) GetProductByName(name string) ([]*entity.Product, error) {
	var returnedProducts []*entity.Product
	if err := r.db.
		Where("name ILIKE ?", "%"+name+"%").
		Find(&returnedProducts).Error; err != nil {
		return nil, err
	}
	return returnedProducts, nil
}

func (r *postgresRepository) DeleteProductByAdmin(id uuid.UUID) error {
	return r.db.Delete(&entity.Product{}, "id = ?", id).Error
}
