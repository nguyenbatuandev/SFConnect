package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}
type ServerConfig struct {
	GinMode string
	Port    string
}

type JWTConfig struct {
	SecretKey string
}

type CallServiceConfig struct {
	UserServiceURL    string
	ProductServiceURL string
	CommissionRate    float64
}

type Config struct {
	Database    DatabaseConfig
	Server      ServerConfig
	JWT         JWTConfig
	CallService CallServiceConfig
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	commissionRateStr := os.Getenv("SERVICE_COMMISSION_RATE")
	commissionRate, err := strconv.ParseFloat(commissionRateStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVICE_COMMISSION_RATE: %v", err)
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:         os.Getenv("DB_HOST"),
			Port:         os.Getenv("DB_PORT"),
			Username:     os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASSWORD"),
			DatabaseName: os.Getenv("DB_NAME"),
			SSLMode:      os.Getenv("DB_SSLMODE"),
		},
		Server: ServerConfig{
			Port:    os.Getenv("SERVER_PORT"),
			GinMode: os.Getenv("GIN_MODE"),
		},
		JWT: JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET"),
		}, CallService: CallServiceConfig{
			UserServiceURL:    os.Getenv("SERVICE_USER_URL"),
			ProductServiceURL: os.Getenv("SERVICE_PRODUCT_URL"),
			CommissionRate:    commissionRate,
		},
	}
	return config, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.Host, c.Username, c.Password, c.DatabaseName, c.Port, c.SSLMode)
}
