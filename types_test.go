package gtmlp

import (
	"testing"
	"time"
)

// TestDefaultEnvMapping verifies the default environment variable names
func TestDefaultEnvMapping(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"Timeout", "GTMLP_TIMEOUT", DefaultEnvMapping.Timeout},
		{"UserAgent", "GTMLP_USER_AGENT", DefaultEnvMapping.UserAgent},
		{"RandomUA", "GTMLP_RANDOM_UA", DefaultEnvMapping.RandomUA},
		{"MaxRetries", "GTMLP_MAX_RETRIES", DefaultEnvMapping.MaxRetries},
		{"Proxy", "GTMLP_PROXY", DefaultEnvMapping.Proxy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.field)
			}
		})
	}

	// Ensure all fields have values
	if DefaultEnvMapping.Timeout == "" {
		t.Error("Timeout should not be empty")
	}
	if DefaultEnvMapping.UserAgent == "" {
		t.Error("UserAgent should not be empty")
	}
	if DefaultEnvMapping.RandomUA == "" {
		t.Error("RandomUA should not be empty")
	}
	if DefaultEnvMapping.MaxRetries == "" {
		t.Error("MaxRetries should not be empty")
	}
	if DefaultEnvMapping.Proxy == "" {
		t.Error("Proxy should not be empty")
	}
}

// TestConfigZeroValues verifies Config struct zero values
func TestConfigZeroValues(t *testing.T) {
	var cfg Config

	// Check string fields have zero values
	if cfg.Container != "" {
		t.Errorf("Expected Container to be empty string, got %q", cfg.Container)
	}
	if cfg.UserAgent != "" {
		t.Errorf("Expected UserAgent to be empty string, got %q", cfg.UserAgent)
	}
	if cfg.Proxy != "" {
		t.Errorf("Expected Proxy to be empty string, got %q", cfg.Proxy)
	}

	// Check map fields are nil
	if cfg.Fields != nil {
		t.Errorf("Expected Fields to be nil, got %v", cfg.Fields)
	}
	if cfg.Headers != nil {
		t.Errorf("Expected Headers to be nil, got %v", cfg.Headers)
	}

	// Check numeric fields have zero values
	if cfg.Timeout != 0 {
		t.Errorf("Expected Timeout to be 0, got %v", cfg.Timeout)
	}
	if cfg.MaxRetries != 0 {
		t.Errorf("Expected MaxRetries to be 0, got %d", cfg.MaxRetries)
	}

	// Check boolean field has zero value
	if cfg.RandomUA {
		t.Error("Expected RandomUA to be false")
	}
}

// TestConfigFieldAssignment verifies Config fields can be properly assigned
func TestConfigFieldAssignment(t *testing.T) {
	cfg := Config{
		Container: "//div[@class='item']",
		Fields: map[string]FieldConfig{
			"title":       {XPath: ".//h2/text()"},
			"description": {XPath: ".//p/text()"},
		},
		Timeout:    30 * time.Second,
		UserAgent:  "TestAgent/1.0",
		RandomUA:   true,
		MaxRetries: 3,
		Proxy:      "http://proxy.example.com:8080",
		Headers: map[string]string{
			"Accept": "text/html",
		},
	}

	if cfg.Container != "//div[@class='item']" {
		t.Errorf("Container assignment failed")
	}

	if len(cfg.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(cfg.Fields))
	}

	if cfg.Fields["title"].XPath != ".//h2/text()" {
		t.Errorf("Fields[\"title\"] assignment failed")
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout assignment failed")
	}

	if !cfg.RandomUA {
		t.Errorf("RandomUA assignment failed")
	}

	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries assignment failed")
	}
}

// TestPartialResultGenericType verifies PartialResult works with different types
func TestPartialResultGenericType(t *testing.T) {
	t.Run("with struct type", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		result := PartialResult[Person]{
			Data: []Person{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 25},
			},
			Errors: map[string]error{
				"name_field": nil,
			},
		}

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 data items, got %d", len(result.Data))
		}

		if result.Data[0].Name != "Alice" {
			t.Errorf("Expected first person to be Alice, got %s", result.Data[0].Name)
		}

		if len(result.Errors) != 1 {
			t.Errorf("Expected 1 error entry, got %d", len(result.Errors))
		}
	})

	t.Run("with primitive string type", func(t *testing.T) {
		result := PartialResult[string]{
			Data: []string{"apple", "banana", "cherry"},
			Errors: map[string]error{
				"parsing": nil,
			},
		}

		if len(result.Data) != 3 {
			t.Errorf("Expected 3 data items, got %d", len(result.Data))
		}

		if result.Data[1] != "banana" {
			t.Errorf("Expected second item to be banana, got %s", result.Data[1])
		}
	})

	t.Run("with primitive int type", func(t *testing.T) {
		result := PartialResult[int]{
			Data:   []int{1, 2, 3, 4, 5},
			Errors: map[string]error{},
		}

		if len(result.Data) != 5 {
			t.Errorf("Expected 5 data items, got %d", len(result.Data))
		}

		if result.Data[2] != 3 {
			t.Errorf("Expected third item to be 3, got %d", result.Data[2])
		}

		if len(result.Errors) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(result.Errors))
		}
	})

	t.Run("with pointer type", func(t *testing.T) {
		type Item struct {
			Value string
		}

		result := PartialResult[*Item]{
			Data: []*Item{
				{Value: "first"},
				{Value: "second"},
			},
			Errors: nil,
		}

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 data items, got %d", len(result.Data))
		}

		if result.Data[0].Value != "first" {
			t.Errorf("Expected first item value to be 'first', got %s", result.Data[0].Value)
		}

		if result.Errors != nil {
			t.Errorf("Expected Errors to be nil, got %v", result.Errors)
		}
	})

	t.Run("with interface type", func(t *testing.T) {
		result := PartialResult[any]{
			Data: []any{
				"string",
				42,
				3.14,
				true,
			},
			Errors: map[string]error{
				"mixed_types": nil,
			},
		}

		if len(result.Data) != 4 {
			t.Errorf("Expected 4 data items, got %d", len(result.Data))
		}

		if result.Data[1] != 42 {
			t.Errorf("Expected second item to be 42, got %v", result.Data[1])
		}
	})
}

// TestPartialResultEmpty verifies PartialResult with empty slices/maps
func TestPartialResultEmpty(t *testing.T) {
	t.Run("empty data and errors", func(t *testing.T) {
		result := PartialResult[string]{
			Data:   []string{},
			Errors: map[string]error{},
		}

		if len(result.Data) != 0 {
			t.Errorf("Expected 0 data items, got %d", len(result.Data))
		}

		if len(result.Errors) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(result.Errors))
		}
	})

	t.Run("nil data and errors", func(t *testing.T) {
		result := PartialResult[int]{
			Data:   nil,
			Errors: nil,
		}

		if result.Data != nil {
			t.Errorf("Expected Data to be nil, got %v", result.Data)
		}

		if result.Errors != nil {
			t.Errorf("Expected Errors to be nil, got %v", result.Errors)
		}
	})
}
