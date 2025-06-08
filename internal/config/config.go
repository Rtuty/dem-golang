package config

import (
	"fmt"
	"os"
)

// Config содержит конфигурацию приложения
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

// ServerConfig содержит конфигурацию сервера
type ServerConfig struct {
	Host string `json:"host" default:"localhost"`
	Port string `json:"port" default:"8080"`
}

// DatabaseConfig содержит конфигурацию базы данных
type DatabaseConfig struct {
	Host     string `json:"host" default:"localhost"`
	Port     string `json:"port" default:"5432"`
	User     string `json:"user" default:"postgres"`
	Password string `json:"password"`
	DBName   string `json:"dbname" default:"wallpaper_system"`
	SSLMode  string `json:"sslmode" default:"disable"`
}

// Load загружает конфигурацию из переменных окружения с дефолтными значениями
func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "wallpaper_system"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	return config
}

// GetDSN возвращает строку подключения к базе данных
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetServerAddress возвращает адрес сервера
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// getEnv получает значение переменной окружения или возвращает дефолтное значение
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
