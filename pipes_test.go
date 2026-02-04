package gtmlp

import (
	"context"
	"testing"
	"time"
)

func TestTrimPipe(t *testing.T) {
	result, err := trimPipe(context.Background(), "  hello  ", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello" {
		t.Errorf("expected 'hello', got '%v'", result)
	}
}

func TestToIntPipe(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		params  []string
		want    any
		wantErr bool
	}{
		{
			name:  "valid integer",
			input: "123",
			want:  123,
		},
		{
			name:  "integer with comma",
			input: "1,234",
			want:  1234,
		},
		{
			name:    "invalid integer",
			input:   "abc",
			wantErr: true,
		},
		{
			name:  "integer with dollar",
			input: "$123",
			want:  123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toIntPipe(context.Background(), tt.input, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("toIntPipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.want {
				t.Errorf("toIntPipe() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestToFloatPipe(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		params  []string
		want    any
		wantErr bool
	}{
		{
			name:  "valid float",
			input: "123.45",
			want:  123.45,
		},
		{
			name:  "float with comma",
			input: "1,234.56",
			want:  1234.56,
		},
		{
			name:  "float with dollar",
			input: "$123.45",
			want:  123.45,
		},
		{
			name:    "invalid float",
			input:   "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toFloatPipe(context.Background(), tt.input, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("toFloatPipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.want {
				t.Errorf("toFloatPipe() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestParseUrlPipe(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "relative URL",
			baseURL: "https://example.com/products",
			input:   "/item/123",
			want:    "https://example.com/item/123",
		},
		{
			name:    "absolute URL",
			baseURL: "https://example.com",
			input:   "https://other.com/page",
			want:    "https://other.com/page",
		},
		{
			name:    "relative path",
			baseURL: "https://example.com/products/page",
			input:   "../items",
			want:    "https://example.com/items",
		},
		{
			name:    "no base URL in context",
			baseURL: "",
			input:   "/item",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.baseURL != "" {
				ctx = WithURL(ctx, tt.baseURL)
			}

			result, err := parseUrlPipe(ctx, tt.input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUrlPipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.want {
				t.Errorf("parseUrlPipe() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestParseTimePipe(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		params  []string
		wantErr bool
	}{
		{
			name:   "valid ISO time",
			input:  "2024-01-15T10:30:00Z",
			params: []string{"2006-01-02T15:04:05Z", "UTC"},
		},
		{
			name:   "missing layout parameter",
			input:  "2024-01-15",
			params: nil,
			wantErr: true,
		},
		{
			name:   "invalid time format",
			input:  "invalid",
			params: []string{"2006-01-02", "UTC"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTimePipe(context.Background(), tt.input, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTimePipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if _, ok := result.(time.Time); !ok {
					t.Errorf("parseTimePipe() should return time.Time, got %T", result)
				}
			}
		})
	}
}

func TestRegexReplacePipe(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		params  []string
		want    string
		wantErr bool
	}{
		{
			name:   "replace spaces",
			input:  "Hello   World",
			params: []string{"\\s+", " "},
			want:   "Hello World",
		},
		{
			name:   "replace with flags",
			input:  "Hello WORLD",
			params: []string{"world", "there", "i"},
			want:   "Hello there",
		},
		{
			name:    "missing parameters",
			input:   "Hello",
			params:  nil,
			wantErr: true,
		},
		{
			name:    "invalid regex",
			input:   "Hello",
			params:  []string{"[invalid", "X"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := regexReplacePipe(context.Background(), tt.input, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("regexReplacePipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.want {
				t.Errorf("regexReplacePipe() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestHumanDurationPipe(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		params  []string
		want    string
		wantErr bool
	}{
		{
			name:  "seconds",
			input: "30",
			want:  "30 seconds ago",
		},
		{
			name:  "minutes",
			input: "120",
			want:  "2 minutes ago",
		},
		{
			name:  "one minute",
			input: "60",
			want:  "1 minute ago",
		},
		{
			name:  "hours",
			input: "7200",
			want:  "2 hours ago",
		},
		{
			name:  "days",
			input: "172800",
			want:  "2 days ago",
		},
		{
			name:    "invalid number",
			input:   "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := humanDurationPipe(context.Background(), tt.input, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("humanDurationPipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.want {
				t.Errorf("humanDurationPipe() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestParsePipeDefinition(t *testing.T) {
	tests := []struct {
		name          string
		definition    string
		expectedName  string
		expectedParams []string
	}{
		{
			name:          "pipe without params",
			definition:    "trim",
			expectedName:  "trim",
			expectedParams: nil,
		},
		{
			name:          "pipe with one param",
			definition:    "parseTime:2006-01-02",
			expectedName:  "parsetime",
			expectedParams: []string{"2006-01-02"},
		},
		{
			name:          "pipe with multiple params",
			definition:    "regexReplace:\\s+: _:i",
			expectedName:  "regexreplace",
			expectedParams: []string{"\\s+", " _", "i"},
		},
		{
			name:          "case insensitive name",
			definition:    "ToUpper",
			expectedName:  "toupper",
			expectedParams: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, params := parsePipeDefinition(tt.definition)
			if name != tt.expectedName {
				t.Errorf("parsePipeDefinition() name = %v, want %v", name, tt.expectedName)
			}
			if len(params) != len(tt.expectedParams) {
				t.Errorf("parsePipeDefinition() params length = %v, want %v", len(params), len(tt.expectedParams))
				return
			}
			for i := range params {
				if params[i] != tt.expectedParams[i] {
					t.Errorf("parsePipeDefinition() params[%d] = %v, want %v", i, params[i], tt.expectedParams[i])
				}
			}
		})
	}
}

func TestRegisterPipe(t *testing.T) {
	// Register a test pipe
	customPipe := func(ctx context.Context, input string, params []string) (any, error) {
		return "custom: " + input, nil
	}

	RegisterPipe("testpipe", customPipe)

	// Verify it's registered (case-insensitive)
	pipe := getPipe("testpipe")
	if pipe == nil {
		t.Fatal("pipe not found")
	}

	pipe2 := getPipe("TESTPIPE")
	if pipe2 == nil {
		t.Fatal("pipe not found (case-insensitive)")
	}

	// Test the pipe
	result, err := pipe(context.Background(), "hello", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "custom: hello" {
		t.Errorf("expected 'custom: hello', got '%v'", result)
	}
}
