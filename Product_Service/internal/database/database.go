package database

import (
	"Product_Service/internal/entity"
	"Product_Service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


func InitDB(cfg *config.Config) (*gorm.DB, error) {
	gormLogger := logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entity.Product{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
