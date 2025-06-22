package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"user_service/internal/entity"
	"log"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *PostgresRepository) GetUserByID(id uuid.UUID) (*entity.User, error) {                                                                                                                                                                                                                                     
	var user entity.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) UpdateUser(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *PostgresRepository) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&entity.User{}, "id = ?", id).Error
}

func (r *PostgresRepository) GetUserByEmail(email string) (*entity.User, error) {
    log.Printf("🔍 Searching for user with email: %s", email)
    
    var user entity.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        log.Printf("❌ Database error: %v", err)
        return nil, err
    }
    
    log.Printf("✅ Found user: %+v", user)
    return &user, nil
}

func (r *PostgresRepository) ToggleUserLockByAdmin(id uuid.UUID, isActive bool) error {
    return r.db.Model(&entity.User{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *PostgresRepository) GetAllUsersByAdmin() ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}