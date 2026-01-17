package pipe

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RegexRule defines a single regex replacement rule.
type RegexRule struct {
	Pattern string
	Replace string
	Flags   string // "i" for case-insensitive, etc.
}

// RegexPipe applies multiple regex replacements.
type RegexPipe struct {
	Rules []RegexRule
}

// Process applies all regex replacement rules to the input.
func (p *RegexPipe) Process(s string) string {
	result := s
	for _, rule := range p.Rules {
		flags := rule.Flags
		pattern := rule.Pattern

		// Handle flags
		if strings.Contains(flags, "i") {
			pattern = "(?i)" + pattern
		}

		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, rule.Replace)
	}
	return result
}

// NumberNormalizePipe converts numbers like "1.5K" to "1500".
type NumberNormalizePipe struct{}

// Process converts shorthand numbers to full numbers.
func (p *NumberNormalizePipe) Process(s string) string {
	s = strings.TrimSpace(s)

	// Remove commas
	s = strings.ReplaceAll(s, ",", "")

	// Check if it's already a number
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return s
	}

	// Multipliers
	multipliers := map[string]float64{
		"K": 1000,
		"M": 1000000,
		"B": 1000000000,
		"T": 1000000000000,
		"k": 1000,
		"m": 1000000,
		"b": 1000000000,
		"t": 1000000000000,
	}

	// Regex to match number with multiplier
	re := regexp.MustCompile(`^([\d.]+)([KMBTkmbt]?)$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 2 {
		return s
	}

	numStr := matches[1]
	suffix := matches[2]

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return s
	}

	if suffix != "" {
		if mult, ok := multipliers[suffix]; ok {
			num *= mult
		}
	}

	// Return as integer string if it's a whole number
	if num == float64(int64(num)) {
		return fmt.Sprintf("%d", int64(num))
	}

	return fmt.Sprintf("%.2f", num)
}

// URLResolvePipe resolves relative URLs against a base URL.
type URLResolvePipe struct {
	BaseURL string
}

// Process resolves a relative URL against the base URL.
func (p *URLResolvePipe) Process(s string) string {
	if s == "" {
		return s
	}

	// If already absolute, return as is
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}

	base, err := url.Parse(p.BaseURL)
	if err != nil {
		return s
	}

	ref, err := url.Parse(s)
	if err != nil {
		return s
	}

	return base.ResolveReference(ref).String()
}

// ExtractEmailPipe extracts an email address from text.
type ExtractEmailPipe struct{}

// Process extracts the first email address found in the input.
func (p *ExtractEmailPipe) Process(s string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	match := re.FindString(s)
	if match != "" {
		return match
	}
	return s
}

// DateFormatPipe converts dates to timestamps.
type DateFormatPipe struct {
	Format string
}

// Process converts a date string to a Unix timestamp.
func (p *DateFormatPipe) Process(s string) string {
	t, err := time.Parse(p.Format, s)
	if err != nil {
		return s
	}
	return fmt.Sprintf("%d", t.Unix())
}

// SubstringPipe extracts a substring.
type SubstringPipe struct {
	Start int
	End   int // -1 means to the end
}

// Process extracts a substring from the input.
func (p *SubstringPipe) Process(s string) string {
	if len(s) == 0 || p.Start >= len(s) {
		return ""
	}

	end := p.End
	if end < 0 || end > len(s) {
		end = len(s)
	}

	return s[p.Start:end]
}

// SplitPipe splits text by a delimiter and returns the first part.
type SplitPipe struct {
	Delimiter string
	Index     int // which part to return (default 0)
}

// Process splits the input and returns the part at Index.
func (p *SplitPipe) Process(s string) string {
	delimiter := p.Delimiter
	if delimiter == "" {
		delimiter = " "
	}

	parts := strings.Split(s, delimiter)
	if p.Index >= 0 && p.Index < len(parts) {
		return parts[p.Index]
	}
	if len(parts) > 0 {
		return parts[0]
	}
	return s
}

// NewRegexPipe creates a new RegexPipe with the given rules.
func NewRegexPipe(rules []RegexRule) *RegexPipe {
	return &RegexPipe{Rules: rules}
}

// NewNumberNormalizePipe creates a new NumberNormalizePipe.
func NewNumberNormalizePipe() *NumberNormalizePipe {
	return &NumberNormalizePipe{}
}

// NewURLResolvePipe creates a new URLResolvePipe.
func NewURLResolvePipe(baseURL string) *URLResolvePipe {
	return &URLResolvePipe{BaseURL: baseURL}
}

// NewExtractEmailPipe creates a new ExtractEmailPipe.
func NewExtractEmailPipe() *ExtractEmailPipe {
	return &ExtractEmailPipe{}
}

// NewDateFormatPipe creates a new DateFormatPipe.
func NewDateFormatPipe(format string) *DateFormatPipe {
	return &DateFormatPipe{Format: format}
}

// NewSubstringPipe creates a new SubstringPipe.
func NewSubstringPipe(start, end int) *SubstringPipe {
	return &SubstringPipe{Start: start, End: end}
}

// NewSplitPipe creates a new SplitPipe.
func NewSplitPipe(delimiter string, index int) *SplitPipe {
	return &SplitPipe{Delimiter: delimiter, Index: index}
}

// ValidateEmailPipe validates if text is an email address.
type ValidateEmailPipe struct{}

// Process returns the email if valid, otherwise empty string.
func (p *ValidateEmailPipe) Process(s string) string {
	s = strings.TrimSpace(s)
	_, err := mail.ParseAddress(s)
	if err != nil {
		return ""
	}
	return s
}

// ValidateURLPipe validates if text is a URL.
type ValidateURLPipe struct{}

// Process returns the URL if valid, otherwise empty string.
func (p *ValidateURLPipe) Process(s string) string {
	s = strings.TrimSpace(s)
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	// Must have a scheme and host to be considered a valid absolute URL
	if u.Scheme == "" || u.Host == "" {
		return ""
	}
	// Scheme must be http or https
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	return s
}

// NewValidateEmailPipe creates a new ValidateEmailPipe.
func NewValidateEmailPipe() *ValidateEmailPipe {
	return &ValidateEmailPipe{}
}

// NewValidateURLPipe creates a new ValidateURLPipe.
func NewValidateURLPipe() *ValidateURLPipe {
	return &ValidateURLPipe{}
}
