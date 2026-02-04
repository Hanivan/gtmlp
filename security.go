package gtmlp

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// validateURL validates a URL according to the config's security settings
// Returns error if URL is invalid or blocked by SSRF protection
func validateURL(rawURL string, config *Config) error {
	// Parse URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Custom validator takes precedence
	if config.URLValidator != nil {
		if err := config.URLValidator(rawURL); err != nil {
			return err
		}
	}

	// SSRF protection (unless AllowPrivateIPs is enabled)
	if !config.AllowPrivateIPs {
		if err := checkSSRF(u); err != nil {
			return err
		}
	}

	// Warn on HTTP (non-HTTPS) usage
	if strings.ToLower(u.Scheme) == "http" {
		getLogger().Warn("http url used",
			"url", rawURL,
			"recommendation", "use_https",
			"risk", "data_transmitted_in_plaintext")
	}

	return nil
}

// checkSSRF checks if URL points to private/internal network (SSRF protection)
func checkSSRF(u *url.URL) error {
	hostname := u.Hostname()

	// Check for localhost variants
	if hostname == "localhost" || hostname == "127.0.0.1" || hostname == "::1" {
		return fmt.Errorf("SSRF protection: localhost access blocked")
	}

	// Resolve hostname to IP addresses
	ips, err := net.LookupIP(hostname)
	if err != nil {
		// DNS lookup failed, but don't block (could be network issue)
		getLogger().Warn("dns lookup failed",
			"hostname", hostname,
			"error", err.Error())
		return nil
	}

	// Check each resolved IP
	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("SSRF protection: private IP address blocked: %s (resolved from %s)", ip, hostname)
		}
	}

	return nil
}

// isPrivateIP checks if an IP address is in a private/internal range
func isPrivateIP(ip net.IP) bool {
	// IPv4 private ranges
	privateIPv4Ranges := []net.IPNet{
		// 10.0.0.0/8
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		// 172.16.0.0/12
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		// 192.168.0.0/16
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
		// 127.0.0.0/8 (loopback)
		{IP: net.IPv4(127, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		// 169.254.0.0/16 (link-local)
		{IP: net.IPv4(169, 254, 0, 0), Mask: net.CIDRMask(16, 32)},
	}

	// Check IPv4 ranges
	for _, privateNet := range privateIPv4Ranges {
		if privateNet.Contains(ip) {
			return true
		}
	}

	// IPv6 checks
	if ip.To4() == nil && len(ip) == net.IPv6len {
		// ::1 (loopback)
		if ip.IsLoopback() {
			return true
		}

		// fc00::/7 (unique local addresses)
		if ip[0] == 0xfc || ip[0] == 0xfd {
			return true
		}

		// fe80::/10 (link-local)
		if ip[0] == 0xfe && (ip[1]&0xc0) == 0x80 {
			return true
		}
	}

	// Check using net package helpers
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

// defaultURLValidator is a basic URL validator that can be used as a template
// Users can provide their own validator via Config.URLValidator
func defaultURLValidator(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	// Require HTTP or HTTPS
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("only http and https schemes are allowed, got: %s", u.Scheme)
	}

	// Require hostname
	if u.Hostname() == "" {
		return fmt.Errorf("URL must have a hostname")
	}

	return nil
}
