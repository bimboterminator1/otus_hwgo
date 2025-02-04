package config

import "fmt"

// Custom error types for configuration handling.
type Error struct {
	Err      error
	Message  string
	Location string // e.g., "server.listener.port"
}

func (e *Error) Error() string {
	return fmt.Sprintf("config error at %s: %s: %v", e.Location, e.Message, e.Err)
}

func NewConfigError(err error, message, location string) *Error {
	return &Error{
		Err:      err,
		Message:  message,
		Location: location,
	}
}

// Config represents the main configuration structure.
type Config struct {
	Components ComponentsConfig `yaml:"components"`
}

// ComponentsConfig holds all component-specific configurations.
type ComponentsConfig struct {
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
	Storage StorageConfig `yaml:"storage"`
}

type ServerConfig struct {
	Listener ListenerConfig `yaml:"listener"`
}

// ListenerConfig holds server listener configurations.
type ListenerConfig struct {
	Protocol string `yaml:"protocol"`
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
}

// LoggingConfig holds logging specific configurations.
type LoggingConfig struct {
	FilePath string `yaml:"file_path"`
	Level    string `yaml:"level"`
	Type     string `yaml:"type"`
}

// StorageConfig holds storage specific configurations.
type StorageConfig struct {
	Type     string            `yaml:"type"` // memory or postgres
	Address  string            `yaml:"address,omitempty"`
	Username string            `yaml:"username,omitempty"`
	Password string            `yaml:"password,omitempty"`
	Database string            `yaml:"database,omitempty"`
	Pool     StoragePoolConfig `yaml:"pool,omitempty"`
}

// StoragePoolConfig holds database connection pool configurations.
type StoragePoolConfig struct {
	MaxOpenConns    int `yaml:"max_open_conns,omitempty"`
	ConnMaxLifetime int `yaml:"conn_max_lifetime,omitempty"` // in seconds
	MaxIdleConns    int `yaml:"max_idle_conns,omitempty"`
}
