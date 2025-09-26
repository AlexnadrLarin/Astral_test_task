package app

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Admin    AdminConfig
}

type ServerConfig struct {
	Port            int
	Host            string
	ShutdownTimeout int
}

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type AdminConfig struct {
	token string
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	loadEnvVars(config)
	return config, nil
}

func loadEnvVars(config *Config) {
	if envVal := os.Getenv("SERVER_HOST"); envVal != "" {
		config.Server.Host = envVal
	}
	if envVal := os.Getenv("SERVER_PORT"); envVal != "" {
		if port, err := strconv.Atoi(envVal); err == nil {
			config.Server.Port = port
		}
	}
	if envVal := os.Getenv("SERVER_SHUTDOWN_TIMEOUT"); envVal != "" {
		if timeout, err := strconv.Atoi(envVal); err == nil {
			config.Server.ShutdownTimeout = timeout
		}
	}

	if envVal := os.Getenv("POSTGRES_HOST"); envVal != "" {
		config.Postgres.Host = envVal
	}
	if envVal := os.Getenv("POSTGRES_PORT"); envVal != "" {
		if port, err := strconv.Atoi(envVal); err == nil {
			config.Postgres.Port = port
		}
	}
	if envVal := os.Getenv("POSTGRES_USER"); envVal != "" {
		config.Postgres.Username = envVal
	}
	if envVal := os.Getenv("POSTGRES_PASSWORD"); envVal != "" {
		config.Postgres.Password = envVal
	}
	if envVal := os.Getenv("POSTGRES_DB"); envVal != "" {
		config.Postgres.DBName = envVal
	}
	if envVal := os.Getenv("POSTGRES_SSL_MODE"); envVal != "" {
		config.Postgres.SSLMode = envVal
	}

	if envVal := os.Getenv("ADMIN_TOKEN"); envVal != "" {
		config.Admin.token = envVal
	}
}
