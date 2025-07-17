package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      AppConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AppConfig struct {
	Environment string
	LogLevel    string
}

func Load() *Config {
	// 環境変数から設定を読み込み
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "jobboard"),
			SSLMode:  getEnv("DB_SSLMODE", "require"),
		},
		App: AppConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
	}

	// 必須項目の確認（本番環境）
	if cfg.IsProduction() {
		validateProductionConfig(cfg)
	}

	log.Printf("Configuration loaded for environment: %s", cfg.App.Environment)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: invalid integer value for %s: %s", key, value)
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return "host=" + c.Database.Host +
		" port=" + strconv.Itoa(c.Database.Port) +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.Name +
		" sslmode=" + c.Database.SSLMode +
		" TimeZone=Asia/Tokyo"
}

func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// 本番環境での必須設定確認
func validateProductionConfig(cfg *Config) {
	required := map[string]string{
		"DB_HOST":     cfg.Database.Host,
		"DB_USER":     cfg.Database.User,
		"DB_PASSWORD": cfg.Database.Password,
		"DB_NAME":     cfg.Database.Name,
	}

	for key, value := range required {
		if value == "" {
			log.Fatalf("Required environment variable %s is not set", key)
		}
	}
}
