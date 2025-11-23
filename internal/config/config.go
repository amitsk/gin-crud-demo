package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Logger      LoggerConfig
	Environment string
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LoggerConfig struct {
	Level  string
	Format string
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("Environment", "test")
	viper.SetDefault("Server.Port", 8080)
	viper.SetDefault("Log.Level", "info")

	// Step 1: Set up environment variable support (highest precedence)
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	env := viper.GetString("ENVIRONMENT")
	if env == "" {
		env = "test" // fallback to default
	}

	// Step 2: Load YAML config (lowest precedence)
	viper.SetConfigName("config." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig() // ignore error if YAML doesn't exist

	// Step 3: Read .env file and manually bind values (middle precedence)
	// Viper doesn't natively support .env files well, so we use a separate viper instance
	envViper := viper.New()
	envViper.SetConfigFile(".env")
	envViper.SetConfigType("env")
	if err := envViper.ReadInConfig(); err == nil {
		// Manually copy values from .env to main viper instance
		// These will override YAML values
		for key, envKey := range map[string]string{
			"Server.Port":       "SERVER_PORT",
			"Server.Mode":       "SERVER_MODE",
			"Database.Host":     "DATABASE_HOST",
			"Database.Port":     "DATABASE_PORT",
			"Database.User":     "DATABASE_USER",
			"Database.Password": "DATABASE_PASSWORD",
			"Database.DBName":   "DATABASE_DBNAME",
			"Database.SSLMode":  "DATABASE_SSLMODE",
			"Logger.Level":      "LOG_LEVEL",
			"Logger.Format":     "LOG_FORMAT",
			"Environment":       "ENVIRONMENT",
		} {
			if val := envViper.Get(envKey); val != nil && val != "" {
				viper.Set(key, val)
			}
		}
	}

	// Step 4: Apply environment variables (highest precedence)
	// Check for APP_ prefixed environment variables and override if present
	bindEnvVars(map[string]string{
		"Server.Port":       "SERVER_PORT",
		"Server.Mode":       "SERVER_MODE",
		"Database.Host":     "DATABASE_HOST",
		"Database.Port":     "DATABASE_PORT",
		"Database.User":     "DATABASE_USER",
		"Database.Password": "DATABASE_PASSWORD",
		"Database.DBName":   "DATABASE_DBNAME",
		"Database.SSLMode":  "DATABASE_SSLMODE",
		"Logger.Level":      "LOG_LEVEL",
		"Logger.Format":     "LOG_FORMAT",
		"Environment":       "ENVIRONMENT",
	})

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func bindEnvVars(mappings map[string]string) {
	for target, key := range mappings {
		// Check for APP_ prefix (from actual environment variables)
		if val := os.Getenv("APP_" + key); val != "" {
			viper.Set(target, val)
		}
	}
}
