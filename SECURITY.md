# Security Policy

## Reporting Security Issues

If you discover a security vulnerability in GTMLP, please email the maintainers privately. Do not create a public GitHub issue.

**Please include:**
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if available)

## Security Features

### SSRF Protection

GTMLP includes built-in Server-Side Request Forgery (SSRF) protection to prevent malicious URL injection attacks.

**What is SSRF?**

SSRF vulnerabilities allow attackers to make requests to internal systems by providing malicious URLs. Example attack vectors:
- `http://localhost:8080/admin` - Access local services
- `http://169.254.169.254/latest/meta-data/` - AWS metadata service
- `http://192.168.1.1/` - Internal network devices

**Protection Enabled by Default:**

GTMLP blocks requests to private IP ranges by default:

```go
config := &gtmlp.Config{
    Container:       "//div",
    Fields:          fields,
    AllowPrivateIPs: false, // default: blocks private IPs
}
```

**Blocked IP Ranges:**
- **Localhost**: `127.0.0.0/8`, `::1`
- **Private IPv4**: `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`
- **Link-local**: `169.254.0.0/16` (AWS/GCP/Azure metadata services)

### URL Scheme Validation

Only HTTP and HTTPS schemes are allowed. Other schemes are rejected:

```go
// ❌ Blocked schemes
"file:///etc/passwd"
"ftp://internal-server"
"javascript:alert(1)"
"data:text/html,<script>alert(1)</script>"
```

### Custom URL Validation

Add domain allowlists or other custom validation:

```go
config := &gtmlp.Config{
    Container: "//div",
    Fields:    fields,

    URLValidator: func(url string) error {
        // Only allow specific domains
        allowedDomains := []string{"example.com", "api.example.com"}

        u, _ := url.Parse(url)
        for _, domain := range allowedDomains {
            if u.Host == domain || strings.HasSuffix(u.Host, "."+domain) {
                return nil
            }
        }

        return fmt.Errorf("domain not in allowlist: %s", u.Host)
    },
}
```

## Security Best Practices

### 1. Never Scrape Untrusted User Input Without Validation

**❌ Insecure:**
```go
userURL := r.URL.Query().Get("url") // User-provided URL
products, _ := gtmlp.ScrapeURL[Product](ctx, userURL, config)
```

**✅ Secure:**
```go
userURL := r.URL.Query().Get("url")

config := &gtmlp.Config{
    Container: "//div",
    Fields:    fields,

    // Validate user input
    URLValidator: func(url string) error {
        u, err := url.Parse(url)
        if err != nil {
            return err
        }

        // Only allow specific domains
        if u.Host != "safe-domain.com" {
            return errors.New("domain not allowed")
        }

        // Only allow HTTPS
        if u.Scheme != "https" {
            return errors.New("only HTTPS allowed")
        }

        return nil
    },
}

products, err := gtmlp.ScrapeURL[Product](ctx, userURL, config)
if err != nil {
    if strings.Contains(err.Error(), "SSRF protection") {
        log.Warn("SSRF attempt blocked", "url", userURL)
    }
    return err
}
```

### 2. Use HTTPS in Production

```go
config := &gtmlp.Config{
    URLValidator: func(url string) error {
        if strings.HasPrefix(url, "http://") {
            return errors.New("HTTP not allowed in production")
        }
        return nil
    },
}
```

### 3. Set Timeouts

Prevent resource exhaustion attacks:

```go
config := &gtmlp.Config{
    Timeout: 30 * time.Second, // Prevent hanging requests

    Pagination: &gtmlp.PaginationConfig{
        MaxPages: 50,              // Limit pages
        Timeout:  10 * time.Minute, // Total pagination timeout
    },
}
```

### 4. Limit Request Rate

Implement rate limiting to prevent abuse:

```go
var rateLimiter = rate.NewLimiter(rate.Every(time.Second), 10) // 10 req/sec

config := &gtmlp.Config{
    URLValidator: func(url string) error {
        if !rateLimiter.Allow() {
            return errors.New("rate limit exceeded")
        }
        return nil
    },
}
```

### 5. Monitor for SSRF Attempts

Log and alert on SSRF attempts:

```go
products, err := gtmlp.ScrapeURL[Product](ctx, url, config)
if err != nil {
    if strings.Contains(err.Error(), "SSRF protection") {
        // Log security event
        securityLog.Warn("SSRF attempt detected",
            "url", url,
            "client_ip", clientIP,
            "user_id", userID)

        // Increment metrics
        metrics.IncrementSSRFAttempts()

        // Alert if threshold exceeded
        if metrics.GetSSRFAttempts() > 100 {
            alerting.TriggerAlert("High SSRF attempt rate")
        }
    }
}
```

### 6. Enable Logging in Production

Monitor scraping activity:

```go
import "log/slog"

// Production: structured JSON logs
handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelWarn, // Log warnings and errors
    AddSource: true,
})
gtmlp.SetLogger(slog.New(handler))
```

### 7. Validate XPath Expressions

Prevent XPath injection:

```go
// ❌ Insecure: user-provided XPath
userXPath := r.URL.Query().Get("xpath")
config.Fields["data"] = gtmlp.FieldConfig{XPath: userXPath}

// ✅ Secure: use predefined XPaths only
predefinedXPaths := map[string]string{
    "product_name":  ".//h2[@class='title']/text()",
    "product_price": ".//span[@class='price']/text()",
}

selectedField := r.URL.Query().Get("field")
if xpath, ok := predefinedXPaths[selectedField]; ok {
    config.Fields["data"] = gtmlp.FieldConfig{XPath: xpath}
} else {
    return errors.New("invalid field selection")
}
```

### 8. Use AllowPrivateIPs Only When Necessary

```go
// ✅ Production: keep SSRF protection enabled
config := &gtmlp.Config{
    AllowPrivateIPs: false, // default
}

// ⚠️ Development/Testing only: allow localhost
if os.Getenv("ENV") == "development" {
    config.AllowPrivateIPs = true
}
```

## Common Attack Vectors

### 1. SSRF via URL Parameter

**Attack:**
```
GET /scrape?url=http://169.254.169.254/latest/meta-data/
```

**Protection:**
```go
// GTMLP blocks this automatically
config.AllowPrivateIPs = false // default
```

### 2. SSRF via Pagination

**Attack:**
```html
<!-- Malicious pagination link -->
<a rel="next" href="http://localhost:8080/admin">Next</a>
```

**Protection:**
```go
// SSRF protection applies to pagination URLs too
config := &gtmlp.Config{
    AllowPrivateIPs: false, // Blocks malicious pagination
    Pagination: &gtmlp.PaginationConfig{
        NextSelector: "//a[@rel='next']/@href",
    },
}
```

### 3. HTTP Scheme Downgrade

**Attack:**
```
GET /scrape?url=http://secure-site.com/api/key
```

**Protection:**
```go
config := &gtmlp.Config{
    URLValidator: func(url string) error {
        if !strings.HasPrefix(url, "https://") {
            return errors.New("HTTPS required")
        }
        return nil
    },
}
```

### 4. DNS Rebinding

**Attack:**
```
evil.com resolves to:
  First:  1.2.3.4 (public IP)
  Second: 127.0.0.1 (localhost)
```

**Protection:**
```go
// GTMLP checks IP after DNS resolution
// If IP is private, request is blocked even if hostname is public
```

### 5. Domain Whitelist Bypass

**Attack:**
```
GET /scrape?url=https://trusted.com.evil.com
```

**Protection:**
```go
config := &gtmlp.Config{
    URLValidator: func(url string) error {
        u, _ := url.Parse(url)

        // ❌ Vulnerable: simple string contains check
        // if strings.Contains(u.Host, "trusted.com") { ... }

        // ✅ Secure: exact or suffix match
        if u.Host == "trusted.com" || strings.HasSuffix(u.Host, ".trusted.com") {
            return nil
        }

        return errors.New("domain not allowed")
    },
}
```

## Security Checklist

Before deploying GTMLP in production:

- [ ] **SSRF Protection**: Keep `AllowPrivateIPs: false` (default)
- [ ] **URL Validation**: Implement `URLValidator` for domain allowlists
- [ ] **HTTPS Only**: Reject HTTP URLs in production
- [ ] **Timeouts**: Set reasonable `Timeout` and pagination limits
- [ ] **Rate Limiting**: Implement per-user/IP rate limits
- [ ] **Input Validation**: Never use user-provided XPath expressions
- [ ] **Logging**: Enable Warn-level logging to monitor activity
- [ ] **Monitoring**: Alert on SSRF attempts and rate limit violations
- [ ] **Error Handling**: Don't expose internal error details to users
- [ ] **Testing**: Test with malicious URLs in staging environment

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 2.x     | :white_check_mark: |
| 1.x     | :warning: Security updates only |
| < 1.0   | :x:                |

## Security Updates

Security updates are released as patch versions (e.g., 2.0.1, 2.0.2).

To stay secure:
```bash
# Update to latest patch version
go get -u github.com/Hanivan/gtmlp@latest

# Check for security advisories
go list -m -u all | grep gtmlp
```

## Additional Resources

- [OWASP SSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Server_Side_Request_Forgery_Prevention_Cheat_Sheet.html)
- [CWE-918: Server-Side Request Forgery (SSRF)](https://cwe.mitre.org/data/definitions/918.html)
- [GTMLP API Documentation](docs/API_V2.md#security)

## License

This security policy is part of the GTMLP project and is licensed under MIT License.
