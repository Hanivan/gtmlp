package gtmlp

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigFormat specifies file format
type ConfigFormat string

const (
	FormatJSON ConfigFormat = "json"
	FormatYAML ConfigFormat = "yaml"
)

// LoadConfig loads selector config from file (JSON/YAML auto-detected)
func LoadConfig(path string, envMapping *EnvMapping) (*Config, error) {
	// Use default mapping if nil
	if envMapping == nil {
		envMapping = DefaultEnvMapping
	}

	// Detect format from extension
	var format ConfigFormat
	if strings.HasSuffix(path, ".json") {
		format = FormatJSON
	} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		format = FormatYAML
	} else {
		return nil, &ScrapeError{
			Type:    ErrTypeConfig,
			Message: fmt.Sprintf("unsupported file format: %s", path),
		}
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ScrapeError{
			Type:    ErrTypeConfig,
			Message: "failed to read config file",
			Cause:   err,
		}
	}

	return ParseConfig(string(data), format, envMapping)
}

// ParseConfig parses config from string
func ParseConfig(data string, format ConfigFormat, envMapping *EnvMapping) (*Config, error) {
	cfg := &Config{
		Fields:     make(map[string]FieldConfig),
		Timeout:    30 * time.Second,
		UserAgent:  "GTMLP/2.0",
		RandomUA:   false,
		MaxRetries: 0,
	}

	var err error

	// Parse based on format
	switch format {
	case FormatJSON:
		err = json.Unmarshal([]byte(data), cfg)
	case FormatYAML:
		err = yaml.Unmarshal([]byte(data), cfg)
	default:
		return nil, &ScrapeError{
			Type:    ErrTypeConfig,
			Message: fmt.Sprintf("unknown format: %s", format),
		}
	}

	if err != nil {
		return nil, &ScrapeError{
			Type:    ErrTypeConfig,
			Message: "failed to parse config",
			Cause:   err,
		}
	}

	// Use default mapping if nil
	if envMapping == nil {
		envMapping = DefaultEnvMapping
	}

	// Apply env variables
	applyEnvVars(cfg, envMapping)

	return cfg, nil
}

// applyEnvVars applies environment variables to config
func applyEnvVars(cfg *Config, mapping *EnvMapping) {
	if v := os.Getenv(mapping.Timeout); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Timeout = d
		}
	}

	if v := os.Getenv(mapping.UserAgent); v != "" {
		cfg.UserAgent = v
	}

	if v := os.Getenv(mapping.RandomUA); v != "" {
		cfg.RandomUA = strings.ToLower(v) == "true" || v == "1"
	}

	if v := os.Getenv(mapping.MaxRetries); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.MaxRetries = n
		}
	}

	if v := os.Getenv(mapping.Proxy); v != "" {
		cfg.Proxy = v
	}
}

// Validate validates the config
func (c *Config) Validate() error {
	if c.Container == "" {
		return &ScrapeError{
			Type:    ErrTypeConfig,
			Message: "container xpath is required",
		}
	}

	if len(c.Fields) == 0 {
		return &ScrapeError{
			Type:    ErrTypeConfig,
			Message: "at least one field is required",
		}
	}

	if c.Timeout <= 0 {
		return &ScrapeError{
			Type:    ErrTypeConfig,
			Message: "timeout must be positive",
		}
	}

	return nil
}
