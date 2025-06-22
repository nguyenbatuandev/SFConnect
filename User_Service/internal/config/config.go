package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig DatabaseConfig
	ServerConfig   ServerConfig
	JwtConfig      JwtConfig
}

type DatabaseConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type JwtConfig struct {
	SecretKey string
	ExpiresIn string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	config := &Config{
		DatabaseConfig: DatabaseConfig{
			Host:         os.Getenv("DB_HOST"),
			Port:         os.Getenv("DB_PORT"),
			Username:     os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASSWORD"),
			DatabaseName: os.Getenv("DB_NAME"),
			SSLMode:      os.Getenv("DB_SSLMODE"),
		},
		ServerConfig: ServerConfig{
			Port:    os.Getenv("SERVER_PORT"),
			GinMode: os.Getenv("GIN_MODE"),
		},
		JwtConfig: JwtConfig{
			SecretKey: os.Getenv("JWT_SECRET"),
			ExpiresIn: os.Getenv("JWT_EXPIRES_IN"),
		},
	}
	return config, nil
}
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.Host, c.Username, c.Password, c.DatabaseName, c.Port, c.SSLMode)
}
