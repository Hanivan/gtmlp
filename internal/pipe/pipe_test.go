package pipe

import (
	"testing"
)

func TestTrimPipe(t *testing.T) {
	pipe := NewTrimPipe()

	result := pipe.Process("  hello world  ")
	if result != "hello world" {
		t.Errorf("TrimPipe failed, got: '%s'", result)
	}
}

func TestLowerCasePipe(t *testing.T) {
	pipe := NewLowerCasePipe()

	result := pipe.Process("HELLO World")
	if result != "hello world" {
		t.Errorf("LowerCasePipe failed, got: '%s'", result)
	}
}

func TestUpperCasePipe(t *testing.T) {
	pipe := NewUpperCasePipe()

	result := pipe.Process("hello World")
	if result != "HELLO WORLD" {
		t.Errorf("UpperCasePipe failed, got: '%s'", result)
	}
}

func TestDecodePipe(t *testing.T) {
	pipe := NewDecodePipe()

	result := pipe.Process("Hello &amp; World")
	if result != "Hello & World" {
		t.Errorf("DecodePipe failed, got: '%s'", result)
	}
}

func TestReplacePipe(t *testing.T) {
	pipe := NewReplacePipe(`\d+`, "NUM")

	result := pipe.Process("Item 1, Item 2, Item 3")
	if result != "Item NUM, Item NUM, Item NUM" {
		t.Errorf("ReplacePipe failed, got: '%s'", result)
	}
}

func TestNumberNormalizePipe(t *testing.T) {
	pipe := NewNumberNormalizePipe()

	tests := []struct {
		input    string
		expected string
	}{
		{"1.5K", "1500"},
		{"2.3M", "2300000"},
		{"500", "500"},
		{"1.2B", "1200000000"},
	}

	for _, tt := range tests {
		result := pipe.Process(tt.input)
		if result != tt.expected {
			t.Errorf("NumberNormalizePipe(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestURLResolvePipe(t *testing.T) {
	pipe := NewURLResolvePipe("https://example.com/path/")

	tests := []struct {
		input    string
		expected string
	}{
		{"/absolute", "https://example.com/absolute"},
		{"relative", "https://example.com/path/relative"},
		{"https://other.com", "https://other.com"},
	}

	for _, tt := range tests {
		result := pipe.Process(tt.input)
		if result != tt.expected {
			t.Errorf("URLResolvePipe(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestExtractEmailPipe(t *testing.T) {
	pipe := NewExtractEmailPipe()

	result := pipe.Process("Contact us at test@example.com for more info")
	if result != "test@example.com" {
		t.Errorf("ExtractEmailPipe failed, got: '%s'", result)
	}

	// Test no email found
	result = pipe.Process("No email here")
	if result != "No email here" {
		t.Errorf("ExtractEmailPipe should return original when no email found, got: '%s'", result)
	}
}

func TestSubstringPipe(t *testing.T) {
	pipe := NewSubstringPipe(0, 5)

	result := pipe.Process("Hello World")
	if result != "Hello" {
		t.Errorf("SubstringPipe failed, got: '%s'", result)
	}

	// Test to end
	pipe2 := NewSubstringPipe(6, -1)
	result = pipe2.Process("Hello World")
	if result != "World" {
		t.Errorf("SubstringPipe(to end) failed, got: '%s'", result)
	}
}

func TestSplitPipe(t *testing.T) {
	pipe := NewSplitPipe(",", 1)

	result := pipe.Process("a,b,c")
	if result != "b" {
		t.Errorf("SplitPipe failed, got: '%s'", result)
	}
}

func TestRegexPipe(t *testing.T) {
	pipe := NewRegexPipe([]RegexRule{
		{Pattern: `red`, Replace: "blue", Flags: "i"},
		{Pattern: `\d+`, Replace: "NUM"},
	})

	result := pipe.Process("I have 5 Red apples and 10 green ones")
	expected := "I have NUM blue apples and NUM green ones"
	if result != expected {
		t.Errorf("RegexPipe failed, got: '%s', expected: '%s'", result, expected)
	}
}

func TestValidateEmailPipe(t *testing.T) {
	pipe := NewValidateEmailPipe()

	valid := pipe.Process("test@example.com")
	if valid != "test@example.com" {
		t.Errorf("ValidateEmailPipe failed for valid email, got: '%s'", valid)
	}

	invalid := pipe.Process("not-an-email")
	if invalid != "" {
		t.Errorf("ValidateEmailPipe should return empty for invalid email, got: '%s'", invalid)
	}
}

func TestValidateURLPipe(t *testing.T) {
	pipe := NewValidateURLPipe()

	valid := pipe.Process("https://example.com")
	if valid != "https://example.com" {
		t.Errorf("ValidateURLPipe failed for valid URL, got: '%s'", valid)
	}

	invalid := pipe.Process("not a url")
	if invalid != "" {
		t.Errorf("ValidateURLPipe should return empty for invalid URL, got: '%s'", invalid)
	}
}

func TestStripHTMLPipe(t *testing.T) {
	pipe := NewStripHTMLPipe()

	result := pipe.Process("<p>Hello <strong>World</strong></p>")
	if result != "Hello World" {
		t.Errorf("StripHTMLPipe failed, got: '%s'", result)
	}
}

func TestPipeRegistry(t *testing.T) {
	// Test getting existing pipe
	pipe, err := CreatePipe("trim")
	if err != nil {
		t.Errorf("CreatePipe failed: %v", err)
	}
	if pipe == nil {
		t.Error("CreatePipe returned nil")
	}

	// Test getting non-existent pipe
	_, err = CreatePipe("nonexistent")
	if err == nil {
		t.Error("CreatePipe should return error for non-existent pipe")
	}

	// Test listing pipes
	pipes := ListPipes()
	if len(pipes) == 0 {
		t.Error("ListPipes returned empty list")
	}

	// Test register/unregister
	RegisterPipe("custom", func() Pipe {
		return NewTrimPipe()
	})

	pipe, err = CreatePipe("custom")
	if err != nil {
		t.Errorf("CreatePipe('custom') failed: %v", err)
	}
	if pipe == nil {
		t.Error("Custom pipe is nil")
	}

	UnregisterPipe("custom")
	_, err = CreatePipe("custom")
	if err == nil {
		t.Error("CreatePipe should return error after unregistering")
	}
}

func TestDateFormatPipe(t *testing.T) {
	pipe := NewDateFormatPipe("2006-01-02")

	result := pipe.Process("2024-01-15")
	// Unix timestamp for 2024-01-15 00:00:00 UTC is 1705276800
	if result != "1705276800" {
		t.Errorf("DateFormatPipe failed, got: '%s'", result)
	}

	// Test invalid date
	result = pipe.Process("invalid-date")
	if result != "invalid-date" {
		t.Errorf("DateFormatPipe should return original for invalid date, got: '%s'", result)
	}
}
