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
	"github.com/thesyscoder/kylon/pkg/logger"
)

var log = logger.GetLogger().WithField("component", "config")

type Config struct {
	App        AppConfig        `yaml:"app"`
	Log        LogConfig        `yaml:"logs"`
	Database   DatabaseConfig   `yaml:"database"`
	Storage    StorageConfig    `yaml:"storage"`
	Kubernetes KubernetesConfig `yaml:"kubernetes"`
	AI         AIConfig         `yaml:"ai"`
	Scheduler  SchedulerConfig  `yaml:"scheduler"`
}

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

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

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

type StorageConfig struct {
	Provider string `yaml:"provider"`
	Bucket   string `yaml:"bucket"`
}

type KubernetesConfig struct {
	KubeconfigPath    string `yaml:"kubeconfigPath"`
	KubeconfigSaveDir string `yaml:"kubeconfigSaveDir"`
}

type AIConfig struct {
	Enabled       bool   `yaml:"enabled"`
	ModelEndpoint string `yaml:"modelEndpoint"`
}

type SchedulerConfig struct {
	IntervalMinutes int `yaml:"intervalMinutes"`
}

var (
	cfg     *Config
	once    sync.Once
	loadErr error
)

func LoadConfig() (*Config, error) {
	once.Do(func() {
		log.Info("Attempting to load application configuration...")

		if envErr := godotenv.Load(); envErr != nil {
			log.Warnf("'.env' file not found or could not be loaded: %v. Proceeding without .env file.", envErr)
		}

		appEnv := os.Getenv("APP_ENV")
		if appEnv == "" {
			appEnv = "development"
			fmt.Fprintf(os.Stdout, "APP_ENV not set. Defaulting to '%s'.\n", appEnv)
		}
		configFilePath := fmt.Sprintf("./configs/%s.config.yaml", appEnv)

		cfg = &Config{}

		if readErr := cleanenv.ReadConfig(configFilePath, cfg); readErr != nil {
			// Removed the specific check for cleanenv.ErrEnvVarUsed
			loadErr = fmt.Errorf("failed to read configuration from '%s': %w", configFilePath, readErr)
			return
		}

		if cfg.App.Env != appEnv {
			log.Warnf("App environment in config file ('%s') differs from determined environment ('%s'). Overriding App.Env with determined value.", cfg.App.Env, appEnv)
			cfg.App.Env = appEnv
		}

		log.WithField("env", cfg.App.Env).WithField("version", cfg.App.Version).Info("Application configuration loaded.")
	})

	return cfg, loadErr
}
