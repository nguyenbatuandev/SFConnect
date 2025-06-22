package config

import (
	"fmt"
	"log"
	"os"

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

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Database int
}
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
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
		}, Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			Database: 0,
		},
	}
	return config, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.Host, c.Username, c.Password, c.DatabaseName, c.Port, c.SSLMode)
}
