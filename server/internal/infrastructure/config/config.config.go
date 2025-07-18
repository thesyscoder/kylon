/**
 * @File: config.go
 * @Title: Application Configuration Management
 * @Description: Handles loading and managing all application configuration settings
 * @Description: from YAML files and environment variables, ensuring a singleton instance.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/thesyscoder/kylon/pkg/logger" // Import the custom logger package
)

// Get the logger instance for this package.
var log = logger.GetLogger().WithField("component", "config")

// Config holds all application configurations, structured into logical sub-sections.
type Config struct {
	App        AppConfig        `yaml:"app"`
	Log        LogConfig        `yaml:"logs"` // Corrected YAML key to "logs"
	Storage    StorageConfig    `yaml:"storage"`
	Kubernetes KubernetesConfig `yaml:"kubernetes"`
	AI         AIConfig         `yaml:"ai"`
	Scheduler  SchedulerConfig  `yaml:"scheduler"`
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Name            string        `yaml:"name"`
	Version         string        `yaml:"version"`
	Host            string        `env:"APP_HOST" yaml:"host"`
	Port            string        `env:"APP_PORT" yaml:"port"`
	Env             string        `env:"APP_ENV" yaml:"env"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"` // e.g., json or text
}

// StorageConfig represents the backup storage provider and bucket info.
type StorageConfig struct {
	Provider string `yaml:"provider"`
	Bucket   string `yaml:"bucket"`
}

// KubernetesConfig holds kubeconfig path for accessing Kubernetes API.
type KubernetesConfig struct {
	KubeconfigPath string `yaml:"kubeconfigPath"`
}

// AIConfig defines settings related to AI-based features.
type AIConfig struct {
	Enabled       bool   `yaml:"enabled"`
	ModelEndpoint string `yaml:"modelEndpoint"`
}

// SchedulerConfig sets up interval-based operations.
type SchedulerConfig struct {
	IntervalMinutes int `yaml:"intervalMinutes"`
}

var (
	cfg     *Config
	once    sync.Once
	loadErr error // Stores any error encountered during the single config load operation
)

// LoadConfig initializes and loads application configuration once.
// It reads from .env files, determines the config file based on APP_ENV,
// and then loads settings from the YAML file, overriding with environment variables.
func LoadConfig() (*Config, error) {
	once.Do(func() {
		// Set logger level to debug early for config loading diagnostics.
		logger.SetLogger("debug")
		log.Info("Initializing configuration...")

		// Load .env variables first if present. Errors are logged but not fatal.
		if envErr := godotenv.Load(); envErr != nil {
			log.Warnf(".env file not found or could not be loaded: %v (proceeding...)", envErr)
		}

		// Determine the application environment from APP_ENV, defaulting to "development".
		appEnv := os.Getenv("APP_ENV")
		if appEnv == "" {
			appEnv = "development"
			log.Infof("APP_ENV not set. Defaulting to '%s'.", appEnv)
		}
		configFilePath := fmt.Sprintf("./configs/%s.config.yaml", appEnv)

		cfg = &Config{} // Initialize the Config struct before reading into it.

		// Read configuration from the determined YAML file path.
		if readErr := cleanenv.ReadConfig(configFilePath, cfg); readErr != nil {
			loadErr = fmt.Errorf("failed to read configuration file '%s': %w", configFilePath, readErr)
			return
		}

		// Ensure environment consistency: prefer the determined APP_ENV over file's.
		if cfg.App.Env != appEnv {
			log.Warnf("APP_ENV in config file (%s) differs from determined env (%s). Overriding with determined value.", cfg.App.Env, appEnv)
			cfg.App.Env = appEnv
		}

		// Set the logger level based on the loaded config's log level.
		logger.SetLogger(cfg.Log.Level)
		log.Info("Configuration loaded successfully.")
	})

	return cfg, loadErr
}
