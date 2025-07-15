/**
 * @File: config.config.go
 * @Title: Application Configuration Management
 * @Description: Handles loading and managing all application configuration settings
 * @Description: from YAML files and environment variables, ensuring a singleton instance.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config holds all application configurations, structured into logical sub-sections.
// This struct serves as the central representation of all application settings.
type Config struct {
	App          AppConfig          `yaml:"app"`
	Log          LogConfig          `yaml:"log"`
	Database     DatabaseConfig     `yaml:"database"`
	Redis        RedisConfig        `yaml:"redis"`
	Auth         AuthConfig         `yaml:"auth"`
	Supabase     SupabaseConfig     `yaml:"supabase"`
	JWT          JWTConfig          `yaml:"jwt"`
	FeatureFlags FeatureFlagsConfig `yaml:"featureFlags"`
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Host            string        `env:"APP_HOST" yaml:"host"`
	Port            string        `env:"APP_PORT" yaml:"port"`
	Env             string        `env:"APP_ENV" yaml:"env"`
	Name            string        `yaml:"name"`
	Version         string        `yaml:"version"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DatabaseConfig holds database connection pool settings.
// [Production Grade] MaxConnections, MaxIdleConnections, ConnectionMaxLifetime are crucial for DB performance.
type DatabaseConfig struct {
	Host                  string        `env:"DB_HOST" yaml:"host"`
	Port                  string        `env:"DB_PORT" yaml:"port"`
	User                  string        `env:"DB_USER" yaml:"user"`
	Password              string        `env:"DB_PASSWORD" yaml:"password"`
	Name                  string        `env:"DB_NAME" yaml:"name"`
	SslMode               string        `env:"DB_SSL_MODE" yaml:"sslMode"`
	MaxConnections        int           `yaml:"maxConnections"`
	MaxIdleConnections    int           `yaml:"maxIdleConnections"`
	ConnectionMaxLifetime time.Duration `yaml:"connectionMaxLifetime"`
}

// RedisConfig holds Redis connection settings.
// [Production Grade] PoolSize, timeouts are important for Redis client performance.
type RedisConfig struct {
	Host         string        `env:"REDIS_HOST" yaml:"host"`
	Port         string        `env:"REDIS_PORT" yaml:"port"`
	Password     string        `env:"REDIS_PASSWORD" yaml:"password"`
	PoolSize     int           `yaml:"poolSize"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

// AuthConfig holds authentication-related settings.
type AuthConfig struct {
	AccessTokenExp  time.Duration `yaml:"accessTokenExp"`
	RefreshTokenExp time.Duration `yaml:"refreshTokenExp"`
}

// SupabaseConfig holds Supabase API keys/URLs.
type SupabaseConfig struct {
	URL            string `env:"SUPABASE_URL" yaml:"url"`
	ServiceRoleKey string `env:"SUPABASE_SERVICE_ROLE_KEY" yaml:"serviceRoleKey"`
	AnonKey        string `env:"SUPABASE_ANON_KEY" yaml:"anonKey"`
}

// JWTConfig holds JWT secret.
type JWTConfig struct {
	Secret string `env:"JWT_SECRET" yaml:"secret"`
}

// FeatureFlagsConfig holds boolean toggles for features.
type FeatureFlagsConfig struct {
	EnableRealtimeUpdates bool `yaml:"enableRealtimeUpdates"`
	EnableAdminDashboard  bool `yaml:"enableAdminDashboard"`
}

var (
	cfg     *Config
	once    sync.Once
	loadErr error
)

// LoadConfig loads the application configuration from the specified YAML file path
// and environment variables. It ensures the configuration is loaded only once.
func LoadConfig() (*Config, error) { // Note: path parameter is removed
	once.Do(func() {
		log.Println("[Config]: Initializing configuration...")

		// Load .env file first. Environment variables from .env will take precedence.
		// [DevOps/SRE] Errors are logged but not fatal, as .env might not exist in all environments.
		if envErr := godotenv.Load(); envErr != nil {
			log.Printf("Warning: .env file not found or could not be loaded: %v", envErr)
			// We don't set loadErr here as missing .env might be expected if env vars are directly set.
		}

		// Determine the environment and construct the config file path
		// [System Design] This allows different configurations for different deployment stages.
		appEnv := os.Getenv("APP_ENV")
		if appEnv == "" {
			appEnv = "development" // Default to development if APP_ENV is not set
			log.Printf("[Config]: APP_ENV not set, defaulting to '%s'.", appEnv)
		}
		configFilePath := fmt.Sprintf("./configs/%s.config.yaml", appEnv)

		cfg = &Config{} // Initialize the Config struct.

		// Read config from the determined YAML file and overlay with environment variables.
		// [Production Grade] `cleanenv` handles precedence: ENV vars (loaded by godotenv) > YAML file.
		if readConfigErr := cleanenv.ReadConfig(configFilePath, cfg); readConfigErr != nil {
			loadErr = fmt.Errorf("[Config]: Failed to read config from %s: %w", configFilePath, readConfigErr)
			return
		}

		// After loading config, verify the APP_ENV in the loaded config matches what was determined.
		// This ensures our struct actually reflects the intended environment.
		if cfg.App.Env != appEnv {
			log.Printf("[Config]: Warning: APP_ENV in config (%s) does not match determined environment (%s). Using determined env for file path.", cfg.App.Env, appEnv)
			cfg.App.Env = appEnv // Ensure the loaded config reflects the environment we loaded for
		}

		log.Println("[Config]: Configuration loaded successfully.")
	})

	return cfg, loadErr
}
