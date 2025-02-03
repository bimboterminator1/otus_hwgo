package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, NewConfigError(err, "failed to open config file", "file_path")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, NewConfigError(err, "failed to read config file", "file_path")
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, NewConfigError(err, "failed to unmarshal config", "yaml")
	}

	// Validate the configuration
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	// Validate server listener port
	if cfg.Components.Server.Listener.Port <= 0 || cfg.Components.Server.Listener.Port > 65535 {
		return NewConfigError(nil, "invalid port number", "server.listener.port")
	}

	// Validate logging level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[cfg.Components.Logging.Level] {
		return NewConfigError(nil, "invalid log level", "server.logging.level")
	}
	validLogForms := map[string]bool{
		"text": true,
		"json": true,
	}
	if !validLogForms[cfg.Components.Logging.Type] {
		return NewConfigError(nil, "invalid log type", "server.logging.type")
	}
	// Validate storage type
	validStorageTypes := map[string]bool{
		"memory":   true,
		"postgres": true,
	}
	if !validStorageTypes[cfg.Components.Storage.Type] {
		return NewConfigError(nil, "invalid storage type", "server.storage.type")
	}

	// Validate PostgreSQL credentials if storage type is postgres
	if cfg.Components.Storage.Type == "postgres" {
		if cfg.Components.Storage.Address == "" {
			return NewConfigError(nil, "address is required for postgres storage", "server.storage.address")
		}
		if cfg.Components.Storage.Username == "" {
			return NewConfigError(nil, "username is required for postgres storage", "server.storage.username")
		}
		if cfg.Components.Storage.Password == "" {
			return NewConfigError(nil, "password is required for postgres storage", "server.storage.password")
		}
	}

	return nil
}

func GetDefaultConfig() *Config {
	return &Config{
		Components: ComponentsConfig{
			Server: ServerConfig{
				Listener: ListenerConfig{
					Port: 8080,
					Host: "0.0.0.0",
				},
			},
			Logging: LoggingConfig{
				FilePath: "/tmp/calendar.log",
				Level:    "info",
			},
			Storage: StorageConfig{
				Type: "memory",
				Pool: StoragePoolConfig{
					MaxOpenConns:    25,
					ConnMaxLifetime: 300,
				},
			},
		},
	}
}
