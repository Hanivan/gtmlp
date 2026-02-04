package gtmlp

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	// Register all built-in pipes
	RegisterPipe("trim", trimPipe)
	RegisterPipe("toint", toIntPipe)
	RegisterPipe("tofloat", toFloatPipe)
	RegisterPipe("parseurl", parseUrlPipe)
	RegisterPipe("parsetime", parseTimePipe)
	RegisterPipe("regexreplace", regexReplacePipe)
	RegisterPipe("humanduration", humanDurationPipe)
}

// trim removes leading/trailing whitespace
func trimPipe(ctx context.Context, input string, params []string) (any, error) {
	return strings.TrimSpace(input), nil
}

// toInt converts string to integer
func toIntPipe(ctx context.Context, input string, params []string) (any, error) {
	// Remove common non-numeric characters
	cleaned := strings.TrimSpace(input)
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")

	val, err := strconv.Atoi(cleaned)
	if err != nil {
		return "", fmt.Errorf("cannot convert '%s' to int: %w", input, err)
	}
	return val, nil
}

// toFloat converts string to float
func toFloatPipe(ctx context.Context, input string, params []string) (any, error) {
	// Remove common non-numeric characters
	cleaned := strings.TrimSpace(input)
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")

	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return "", fmt.Errorf("cannot convert '%s' to float: %w", input, err)
	}
	return val, nil
}

// parseUrl converts relative URLs to absolute using base URL from context
func parseUrlPipe(ctx context.Context, input string, params []string) (any, error) {
	// Get base URL from context
	baseURLVal := ctx.Value(contextKey("baseURL"))
	if baseURLVal == nil {
		return "", fmt.Errorf("baseURL not found in context")
	}

	baseURL, ok := baseURLVal.(string)
	if !ok {
		return "", fmt.Errorf("baseURL is not a string")
	}

	// Parse base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL '%s': %w", baseURL, err)
	}

	// Resolve relative URL
	relative, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("invalid relative URL '%s': %w", input, err)
	}

	absURL := base.ResolveReference(relative)
	return absURL.String(), nil
}

// parseTime parses date/time string with timezone
// Params: [layout, timezone]
// Example: "parsetime:2006-01-02T15:04:05Z:America/New_York"
func parseTimePipe(ctx context.Context, input string, params []string) (any, error) {
	if len(params) < 1 {
		return "", fmt.Errorf("parseTime requires layout parameter (e.g., parseTime:2006-01-02:UTC)")
	}

	layout := params[0]
	timezone := "UTC"
	if len(params) >= 2 {
		timezone = params[1]
	}

	// Parse with timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}

	t, err := time.ParseInLocation(layout, input, loc)
	if err != nil {
		return "", fmt.Errorf("cannot parse time '%s' with layout '%s': %w", input, layout, err)
	}

	return t, nil
}

// regexReplace performs regex substitution
// Params: [pattern, replacement, flags]
// Example: "regexReplace:\\s+:_:i" for case-insensitive
func regexReplacePipe(ctx context.Context, input string, params []string) (any, error) {
	if len(params) < 2 {
		return "", fmt.Errorf("regexReplace requires pattern and replacement (e.g., regexReplace:\\d+:X)")
	}

	pattern := params[0]
	replacement := params[1]
	flags := ""
	if len(params) >= 3 {
		flags = params[2]
	}

	// Build regex with flags
	reStr := pattern
	if strings.Contains(flags, "i") {
		reStr = "(?i)" + pattern
	}

	re, err := regexp.Compile(reStr)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern '%s': %w", pattern, err)
	}

	result := re.ReplaceAllString(input, replacement)
	return result, nil
}

// humanDuration converts seconds to human-readable format
func humanDurationPipe(ctx context.Context, input string, params []string) (any, error) {
	// Try to parse as int first
	seconds, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return "", fmt.Errorf("cannot convert '%s' to seconds: %w", input, err)
	}

	duration := time.Duration(seconds) * time.Second
	return humanizeDuration(duration), nil
}

// humanizeDuration converts duration to human-readable string
func humanizeDuration(d time.Duration) string {
	seconds := int(d.Seconds())

	if seconds < 60 {
		return fmt.Sprintf("%d seconds ago", seconds)
	}

	minutes := seconds / 60
	if minutes < 60 {
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	hours := minutes / 60
	if hours < 24 {
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	days := hours / 24
	if days == 1 {
		return "1 day ago"
	}
	return fmt.Sprintf("%d days ago", days)
}
