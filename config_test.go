package gtmlp

import (
	"os"
	"testing"
	"time"
)

// TestLoadConfig_Success_JSON tests loading a valid JSON config file
func TestLoadConfig_Success_JSON(t *testing.T) {
	cfg, err := LoadConfig("testdata/config.json", nil)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Container != "//div[@class='product']" {
		t.Errorf("expected container '//div[@class='product']', got '%s'", cfg.Container)
	}

	if len(cfg.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(cfg.Fields))
	}

	if cfg.Fields["name"].XPath != ".//h2/text()" {
		t.Errorf("expected field 'name' to be './/h2/text()', got '%s'", cfg.Fields["name"].XPath)
	}

	// Check defaults
	if cfg.Timeout != 30*time.Second {
		t.Errorf("expected default timeout 30s, got %v", cfg.Timeout)
	}

	if cfg.UserAgent != "GTMLP/2.0" {
		t.Errorf("expected default user agent 'GTMLP/2.0', got '%s'", cfg.UserAgent)
	}

	if cfg.RandomUA != false {
		t.Errorf("expected default randomUA false, got %v", cfg.RandomUA)
	}
}

// TestLoadConfig_Success_YAML tests loading a valid YAML config file
func TestLoadConfig_Success_YAML(t *testing.T) {
	cfg, err := LoadConfig("testdata/config.yaml", nil)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Container != "//div[@class='product']" {
		t.Errorf("expected container '//div[@class='product']', got '%s'", cfg.Container)
	}

	if len(cfg.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(cfg.Fields))
	}

	if cfg.Fields["name"].XPath != ".//h2/text()" {
		t.Errorf("expected field 'name' to be './/h2/text()', got '%s'", cfg.Fields["name"].XPath)
	}

	if cfg.Fields["price"].XPath != ".//span[@class='price']/text()" {
		t.Errorf("expected field 'price' to be './/span[@class='price']/text()', got '%s'", cfg.Fields["price"].XPath)
	}
}

// TestLoadConfig_UnsupportedExtension tests error handling for unsupported file formats
func TestLoadConfig_UnsupportedExtension(t *testing.T) {
	_, err := LoadConfig("testdata/config.xml", nil)
	if err == nil {
		t.Fatal("expected error for unsupported file format, got nil")
	}

	if !Is(err, ErrTypeConfig) {
		t.Errorf("expected ErrTypeConfig, got %v", err)
	}

	scrapeErr := err.(*ScrapeError)
	if scrapeErr.Message != "unsupported file format: testdata/config.xml" {
		t.Errorf("expected 'unsupported file format' message, got '%s'", scrapeErr.Message)
	}
}

// TestLoadConfig_FileNotFound tests error handling when file doesn't exist
func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("testdata/nonexistent.json", nil)
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}

	if !Is(err, ErrTypeConfig) {
		t.Errorf("expected ErrTypeConfig, got %v", err)
	}
}

// TestLoadConfig_NilEnvMapping tests that nil envMapping uses defaults
func TestLoadConfig_NilEnvMapping(t *testing.T) {
	// Set an environment variable with default mapping
	os.Setenv("GTMLP_TIMEOUT", "45s")
	defer os.Unsetenv("GTMLP_TIMEOUT")

	cfg, err := LoadConfig("testdata/config.json", nil)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Timeout != 45*time.Second {
		t.Errorf("expected timeout 45s from env var, got %v", cfg.Timeout)
	}
}

// TestLoadConfig_CustomEnvMapping tests custom environment variable mapping
func TestLoadConfig_CustomEnvMapping(t *testing.T) {
	// Set custom env vars
	os.Setenv("CUSTOM_TIMEOUT", "60s")
	os.Setenv("CUSTOM_USER_AGENT", "CustomAgent/1.0")
	defer os.Unsetenv("CUSTOM_TIMEOUT")
	defer os.Unsetenv("CUSTOM_USER_AGENT")

	customMapping := &EnvMapping{
		Timeout:    "CUSTOM_TIMEOUT",
		UserAgent:  "CUSTOM_USER_AGENT",
		RandomUA:   "CUSTOM_RANDOM_UA",
		MaxRetries: "CUSTOM_MAX_RETRIES",
		Proxy:      "CUSTOM_PROXY",
	}

	cfg, err := LoadConfig("testdata/config.json", customMapping)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Timeout != 60*time.Second {
		t.Errorf("expected timeout 60s from custom env var, got %v", cfg.Timeout)
	}

	if cfg.UserAgent != "CustomAgent/1.0" {
		t.Errorf("expected user agent 'CustomAgent/1.0' from custom env var, got '%s'", cfg.UserAgent)
	}
}

// TestParseConfig_JSON tests JSON parsing
func TestParseConfig_JSON(t *testing.T) {
	jsonData := `{
		"container": "//div[@class='item']",
		"fields": {
			"title": {"xpath": ".//h1/text()"},
			"desc": {"xpath": ".//p/text()"}
		}
	}`

	cfg, err := ParseConfig(jsonData, FormatJSON, DefaultEnvMapping)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	if cfg.Container != "//div[@class='item']" {
		t.Errorf("expected container '//div[@class='item']', got '%s'", cfg.Container)
	}

	if len(cfg.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(cfg.Fields))
	}

	if cfg.Fields["title"].XPath != ".//h1/text()" {
		t.Errorf("expected field 'title' to be './/h1/text()', got '%s'", cfg.Fields["title"].XPath)
	}

	if cfg.Fields["desc"].XPath != ".//p/text()" {
		t.Errorf("expected field 'desc' to be './/p/text()', got '%s'", cfg.Fields["desc"].XPath)
	}

	// Check defaults
	if cfg.Timeout != 30*time.Second {
		t.Errorf("expected default timeout 30s, got %v", cfg.Timeout)
	}

	if cfg.UserAgent != "GTMLP/2.0" {
		t.Errorf("expected default user agent 'GTMLP/2.0', got '%s'", cfg.UserAgent)
	}
}

// TestParseConfig_YAML tests YAML parsing
func TestParseConfig_YAML(t *testing.T) {
	yamlData := `container: "//div[@class='item']"
fields:
  title:
    xpath: ".//h1/text()"
  desc:
    xpath: ".//p/text()"
`

	cfg, err := ParseConfig(yamlData, FormatYAML, DefaultEnvMapping)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	if cfg.Container != "//div[@class='item']" {
		t.Errorf("expected container '//div[@class='item']', got '%s'", cfg.Container)
	}

	if len(cfg.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(cfg.Fields))
	}

	if cfg.Fields["title"].XPath != ".//h1/text()" {
		t.Errorf("expected field 'title' to be './/h1/text()', got '%s'", cfg.Fields["title"].XPath)
	}

	if cfg.Fields["desc"].XPath != ".//p/text()" {
		t.Errorf("expected field 'desc' to be './/p/text()', got '%s'", cfg.Fields["desc"].XPath)
	}
}

// TestParseConfig_InvalidJSON tests error handling for invalid JSON
func TestParseConfig_InvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json content`

	_, err := ParseConfig(invalidJSON, FormatJSON, DefaultEnvMapping)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}

	if !Is(err, ErrTypeConfig) {
		t.Errorf("expected ErrTypeConfig, got %v", err)
	}

	scrapeErr := err.(*ScrapeError)
	if scrapeErr.Message != "failed to parse config" {
		t.Errorf("expected 'failed to parse config' message, got '%s'", scrapeErr.Message)
	}

	if scrapeErr.Cause == nil {
		t.Error("expected Cause to be set with parsing error")
	}
}

// TestParseConfig_InvalidYAML tests error handling for invalid YAML
func TestParseConfig_InvalidYAML(t *testing.T) {
	invalidYAML := `
container: "//div"
  invalid_indentation: "bad"
    more_bad: "stuff"
`

	_, err := ParseConfig(invalidYAML, FormatYAML, DefaultEnvMapping)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}

	if !Is(err, ErrTypeConfig) {
		t.Errorf("expected ErrTypeConfig, got %v", err)
	}
}

// TestParseConfig_DefaultValues tests that default values are applied
func TestParseConfig_DefaultValues(t *testing.T) {
	minimalData := `{
		"container": "//div",
		"fields": {"title": {"xpath": ".//h1"}}
	}`

	cfg, err := ParseConfig(minimalData, FormatJSON, DefaultEnvMapping)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("expected default timeout 30s, got %v", cfg.Timeout)
	}

	if cfg.UserAgent != "GTMLP/2.0" {
		t.Errorf("expected default user agent 'GTMLP/2.0', got '%s'", cfg.UserAgent)
	}

	if cfg.RandomUA != false {
		t.Errorf("expected default randomUA false, got %v", cfg.RandomUA)
	}

	if cfg.MaxRetries != 0 {
		t.Errorf("expected default maxRetries 0, got %d", cfg.MaxRetries)
	}
}

// TestApplyEnvVars_Timeout tests timeout environment variable application
func TestApplyEnvVars_Timeout(t *testing.T) {
	os.Setenv("GTMLP_TIMEOUT", "90s")
	defer os.Unsetenv("GTMLP_TIMEOUT")

	cfg := &Config{Timeout: 30 * time.Second}
	applyEnvVars(cfg, DefaultEnvMapping)

	if cfg.Timeout != 90*time.Second {
		t.Errorf("expected timeout 90s from env var, got %v", cfg.Timeout)
	}
}

// TestApplyEnvVars_UserAgent tests user agent environment variable application
func TestApplyEnvVars_UserAgent(t *testing.T) {
	os.Setenv("GTMLP_USER_AGENT", "MyAgent/2.0")
	defer os.Unsetenv("GTMLP_USER_AGENT")

	cfg := &Config{UserAgent: "Default"}
	applyEnvVars(cfg, DefaultEnvMapping)

	if cfg.UserAgent != "MyAgent/2.0" {
		t.Errorf("expected user agent 'MyAgent/2.0' from env var, got '%s'", cfg.UserAgent)
	}
}

// TestApplyEnvVars_RandomUA tests random user agent environment variable application
func TestApplyEnvVars_RandomUA(t *testing.T) {
	testCases := []struct {
		envValue string
		expected bool
		testName string
	}{
		{"true", true, "lowercase true"},
		{"TRUE", true, "uppercase true"},
		{"1", true, "numeric 1"},
		{"false", false, "lowercase false"},
		{"FALSE", false, "uppercase false"},
		{"0", false, "numeric 0"},
		{"invalid", false, "invalid value"},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			os.Setenv("GTMLP_RANDOM_UA", tc.envValue)
			defer os.Unsetenv("GTMLP_RANDOM_UA")

			cfg := &Config{RandomUA: false}
			applyEnvVars(cfg, DefaultEnvMapping)

			if cfg.RandomUA != tc.expected {
				t.Errorf("expected RandomUA %v for env value '%s', got %v", tc.expected, tc.envValue, cfg.RandomUA)
			}
		})
	}
}

// TestApplyEnvVars_MaxRetries tests max retries environment variable application
func TestApplyEnvVars_MaxRetries(t *testing.T) {
	os.Setenv("GTMLP_MAX_RETRIES", "5")
	defer os.Unsetenv("GTMLP_MAX_RETRIES")

	cfg := &Config{MaxRetries: 0}
	applyEnvVars(cfg, DefaultEnvMapping)

	if cfg.MaxRetries != 5 {
		t.Errorf("expected maxRetries 5 from env var, got %d", cfg.MaxRetries)
	}
}

// TestApplyEnvVars_Proxy tests proxy environment variable application
func TestApplyEnvVars_Proxy(t *testing.T) {
	os.Setenv("GTMLP_PROXY", "http://proxy.example.com:8080")
	defer os.Unsetenv("GTMLP_PROXY")

	cfg := &Config{Proxy: ""}
	applyEnvVars(cfg, DefaultEnvMapping)

	if cfg.Proxy != "http://proxy.example.com:8080" {
		t.Errorf("expected proxy 'http://proxy.example.com:8080' from env var, got '%s'", cfg.Proxy)
	}
}

// TestApplyEnvVars_InvalidTimeout tests that invalid timeout values are ignored
func TestApplyEnvVars_InvalidTimeout(t *testing.T) {
	os.Setenv("GTMLP_TIMEOUT", "invalid-duration")
	defer os.Unsetenv("GTMLP_TIMEOUT")

	cfg := &Config{Timeout: 30 * time.Second}
	applyEnvVars(cfg, DefaultEnvMapping)

	if cfg.Timeout != 30*time.Second {
		t.Errorf("expected timeout to remain 30s with invalid env var, got %v", cfg.Timeout)
	}
}

// TestApplyEnvVars_InvalidMaxRetries tests that invalid max retries values are ignored
func TestApplyEnvVars_InvalidMaxRetries(t *testing.T) {
	os.Setenv("GTMLP_MAX_RETRIES", "not-a-number")
	defer os.Unsetenv("GTMLP_MAX_RETRIES")

	cfg := &Config{MaxRetries: 3}
	applyEnvVars(cfg, DefaultEnvMapping)

	if cfg.MaxRetries != 3 {
		t.Errorf("expected maxRetries to remain 3 with invalid env var, got %d", cfg.MaxRetries)
	}
}

// TestValidate_Success tests validation of valid config
func TestValidate_Success(t *testing.T) {
	cfg := &Config{
		Container: "//div[@class='product']",
		Fields: map[string]FieldConfig{
			"name": {XPath: ".//h2/text()"},
		},
		Timeout: 30 * time.Second,
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate failed for valid config: %v", err)
	}
}

// TestValidate_EmptyContainer tests error when container is empty
func TestValidate_EmptyContainer(t *testing.T) {
	cfg := &Config{
		Container: "",
		Fields: map[string]FieldConfig{
			"name": {XPath: ".//h2/text()"},
		},
		Timeout: 30 * time.Second,
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for empty container, got nil")
	}

	if !Is(err, ErrTypeConfig) {
		t.Errorf("expected ErrTypeConfig, got %v", err)
	}

	scrapeErr := err.(*ScrapeError)
	if scrapeErr.Message != "container xpath is required" {
		t.Errorf("expected 'container xpath is required' message, got '%s'", scrapeErr.Message)
	}
}

// TestValidate_EmptyFields tests error when fields map is empty
func TestValidate_EmptyFields(t *testing.T) {
	cfg := &Config{
		Container: "//div[@class='product']",
		Fields:    map[string]FieldConfig{},
		Timeout:   30 * time.Second,
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for empty fields, got nil")
	}

	if !Is(err, ErrTypeConfig) {
		t.Errorf("expected ErrTypeConfig, got %v", err)
	}

	scrapeErr := err.(*ScrapeError)
	if scrapeErr.Message != "at least one field is required" {
		t.Errorf("expected 'at least one field is required' message, got '%s'", scrapeErr.Message)
	}
}

// TestValidate_NilFields tests error when fields map is nil
func TestValidate_NilFields(t *testing.T) {
	cfg := &Config{
		Container: "//div[@class='product']",
		Fields:    nil,
		Timeout:   30 * time.Second,
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for nil fields, got nil")
	}

	if !Is(err, ErrTypeConfig) {
		t.Errorf("expected ErrTypeConfig, got %v", err)
	}
}

// TestValidate_ZeroTimeout tests error when timeout is zero
func TestValidate_ZeroTimeout(t *testing.T) {
	cfg := &Config{
		Container: "//div[@class='product']",
		Fields: map[string]FieldConfig{
			"name": {XPath: ".//h2/text()"},
		},
		Timeout: 0,
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for zero timeout, got nil")
	}

	if !Is(err, ErrTypeConfig) {
		t.Errorf("expected ErrTypeConfig, got %v", err)
	}

	scrapeErr := err.(*ScrapeError)
	if scrapeErr.Message != "timeout must be positive" {
		t.Errorf("expected 'timeout must be positive' message, got '%s'", scrapeErr.Message)
	}
}

// TestValidate_NegativeTimeout tests error when timeout is negative
func TestValidate_NegativeTimeout(t *testing.T) {
	cfg := &Config{
		Container: "//div[@class='product']",
		Fields: map[string]FieldConfig{
			"name": {XPath: ".//h2/text()"},
		},
		Timeout: -10 * time.Second,
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative timeout, got nil")
	}

	if !Is(err, ErrTypeConfig) {
		t.Errorf("expected ErrTypeConfig, got %v", err)
	}

	scrapeErr := err.(*ScrapeError)
	if scrapeErr.Message != "timeout must be positive" {
		t.Errorf("expected 'timeout must be positive' message, got '%s'", scrapeErr.Message)
	}
}

// TestValidate_InvalidAltContainer tests error when altContainer has invalid XPath
func TestValidate_InvalidAltContainer(t *testing.T) {
	cfg := &Config{
		Container:    "//div[@class='product']",
		AltContainer: []string{"//valid", "//invalid[[["},
		Fields: map[string]FieldConfig{
			"name": {XPath: ".//h2/text()"},
		},
		Timeout: 30 * time.Second,
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid altContainer XPath, got nil")
	}

	if !Is(err, ErrTypeXPath) {
		t.Errorf("expected ErrTypeXPath, got %v", err)
	}
}

// TestValidate_InvalidAltXPath tests error when altXpath has invalid XPath
func TestValidate_InvalidAltXPath(t *testing.T) {
	cfg := &Config{
		Container: "//div[@class='product']",
		Fields: map[string]FieldConfig{
			"name": {
				XPath:    ".//h2/text()",
				AltXPath: []string{".//h1/text()", ".//invalid[[[@"},
			},
		},
		Timeout: 30 * time.Second,
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid altXpath, got nil")
	}

	if !Is(err, ErrTypeXPath) {
		t.Errorf("expected ErrTypeXPath, got %v", err)
	}
}

// TestValidate_ValidAltXPath tests that valid altXpath passes validation
func TestValidate_ValidAltXPath(t *testing.T) {
	cfg := &Config{
		Container:    "//div[@class='product']",
		AltContainer: []string{"//div[@class='item']", "//article"},
		Fields: map[string]FieldConfig{
			"name": {
				XPath:    ".//h2/text()",
				AltXPath: []string{".//h1/text()", ".//h3/text()"},
			},
			"price": {
				XPath:    ".//span[@class='price']/text()",
				AltXPath: []string{".//div[@class='price']/text()"},
			},
		},
		Timeout: 30 * time.Second,
	}

	err := cfg.Validate()
	if err != nil {
		t.Fatalf("expected valid config with altXpath to pass validation, got error: %v", err)
	}
}

// TestValidate_EmptyAltArrays tests that empty alt arrays are valid
func TestValidate_EmptyAltArrays(t *testing.T) {
	cfg := &Config{
		Container:    "//div[@class='product']",
		AltContainer: []string{},
		Fields: map[string]FieldConfig{
			"name": {
				XPath:    ".//h2/text()",
				AltXPath: []string{},
			},
		},
		Timeout: 30 * time.Second,
	}

	err := cfg.Validate()
	if err != nil {
		t.Fatalf("expected config with empty alt arrays to pass validation, got error: %v", err)
	}
}

// TestParseConfig_WithAltXPath tests parsing JSON config with altXpath
func TestParseConfig_WithAltXPath(t *testing.T) {
	jsonConfig := `{
		"container": "//div[@class='product']",
		"altContainer": ["//div[@class='item']"],
		"fields": {
			"name": {
				"xpath": ".//h2/text()",
				"altXpath": [".//h1/text()", ".//h3/text()"]
			}
		}
	}`

	cfg, err := ParseConfig(jsonConfig, FormatJSON, nil)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	if len(cfg.AltContainer) != 1 {
		t.Errorf("expected 1 altContainer, got %d", len(cfg.AltContainer))
	}

	if cfg.AltContainer[0] != "//div[@class='item']" {
		t.Errorf("expected altContainer '//div[@class='item']', got '%s'", cfg.AltContainer[0])
	}

	if len(cfg.Fields["name"].AltXPath) != 2 {
		t.Errorf("expected 2 altXpath entries, got %d", len(cfg.Fields["name"].AltXPath))
	}

	if cfg.Fields["name"].AltXPath[0] != ".//h1/text()" {
		t.Errorf("expected first altXpath './/h1/text()', got '%s'", cfg.Fields["name"].AltXPath[0])
	}
}
