package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config представляет конфигурацию приложения
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logger   LoggerConfig   `json:"logger"`
	App      AppConfig      `json:"app"`
}

// ServerConfig содержит конфигурацию HTTP сервера
type ServerConfig struct {
	Host            string        `json:"host"`
	Port            string        `json:"port"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
	EnableCORS      bool          `json:"enable_cors"`
	TrustedProxies  []string      `json:"trusted_proxies"`
}

// DatabaseConfig содержит конфигурацию базы данных
type DatabaseConfig struct {
	Host            string        `json:"host"`
	Port            string        `json:"port"`
	User            string        `json:"user"`
	Password        string        `json:"password"`
	DBName          string        `json:"dbname"`
	SSLMode         string        `json:"sslmode"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

// LoggerConfig содержит конфигурацию логирования
type LoggerConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"` // json или text
	Output     string `json:"output"` // stdout, stderr или путь к файлу
	EnableFile bool   `json:"enable_file"`
	FilePath   string `json:"file_path"`
	MaxSize    int    `json:"max_size"`    // в MB
	MaxBackups int    `json:"max_backups"` // количество резервных файлов
	MaxAge     int    `json:"max_age"`     // в днях
	Compress   bool   `json:"compress"`
}

// AppConfig содержит общие настройки приложения
type AppConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"` // development, staging, production
	Debug       bool   `json:"debug"`
	Timezone    string `json:"timezone"`
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	// Пытаемся загрузить .env файл (игнорируем ошибку если файла нет)
	_ = godotenv.Load()

	config := &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		Logger:   loadLoggerConfig(),
		App:      loadAppConfig(),
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	return config, nil
}

// loadServerConfig загружает конфигурацию сервера
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Host:            getEnv("SERVER_HOST", "localhost"),
		Port:            getEnv("SERVER_PORT", "8080"),
		ReadTimeout:     getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout:    getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:     getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		EnableCORS:      getEnvBool("SERVER_ENABLE_CORS", true),
		TrustedProxies:  getEnvSlice("SERVER_TRUSTED_PROXIES", []string{}),
	}
}

// loadDatabaseConfig загружает конфигурацию базы данных
func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "wallpaper_system"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}
}

// loadLoggerConfig загружает конфигурацию логирования
func loadLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:      getEnv("LOG_LEVEL", "info"),
		Format:     getEnv("LOG_FORMAT", "json"),
		Output:     getEnv("LOG_OUTPUT", "stdout"),
		EnableFile: getEnvBool("LOG_ENABLE_FILE", false),
		FilePath:   getEnv("LOG_FILE_PATH", "logs/app.log"),
		MaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
		MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 3),
		MaxAge:     getEnvInt("LOG_MAX_AGE", 28),
		Compress:   getEnvBool("LOG_COMPRESS", true),
	}
}

// loadAppConfig загружает общие настройки приложения
func loadAppConfig() AppConfig {
	return AppConfig{
		Name:        getEnv("APP_NAME", "Wallpaper System"),
		Version:     getEnv("APP_VERSION", "1.0.0"),
		Environment: getEnv("APP_ENV", "development"),
		Debug:       getEnvBool("APP_DEBUG", false),
		Timezone:    getEnv("APP_TIMEZONE", "UTC"),
	}
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("хост базы данных не может быть пустым")
	}

	if c.Database.User == "" {
		return fmt.Errorf("пользователь базы данных не может быть пустым")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("имя базы данных не может быть пустым")
	}

	if c.Server.Port == "" {
		return fmt.Errorf("порт сервера не может быть пустым")
	}

	// Проверяем валидность уровня логирования
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true, "panic": true,
	}
	if !validLogLevels[c.Logger.Level] {
		return fmt.Errorf("неверный уровень логирования: %s", c.Logger.Level)
	}

	// Проверяем валидность формата логирования
	if c.Logger.Format != "json" && c.Logger.Format != "text" {
		return fmt.Errorf("неверный формат логирования: %s", c.Logger.Format)
	}

	// Проверяем валидность окружения
	validEnvs := map[string]bool{
		"development": true, "staging": true, "production": true,
	}
	if !validEnvs[c.App.Environment] {
		return fmt.Errorf("неверное окружение: %s", c.App.Environment)
	}

	return nil
}

// GetDSN возвращает строку подключения к базе данных
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// GetServerAddress возвращает адрес сервера
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// IsDevelopment проверяет, является ли окружение разработческим
func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction проверяет, является ли окружение продакшеном
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsStaging проверяет, является ли окружение тестовым
func (c *AppConfig) IsStaging() bool {
	return c.Environment == "staging"
}

// Вспомогательные функции для загрузки переменных окружения

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	// Для простоты возвращаем значение по умолчанию
	// В реальном проекте здесь можно реализовать парсинг строки с разделителями
	return defaultValue
} 