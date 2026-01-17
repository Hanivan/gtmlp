# Examples

## Basic Examples

The `examples/basic/` directory contains 7 foundational examples:

1. **XPath Selectors** (`01-xpath.go`) - Query HTML with XPath expressions
2. **JSON Conversion** (`02-json.go`) - Convert HTML to structured JSON
3. **Fluent API Builder** (`03-builder.go`) - Chainable API usage
4. **URL Fetching** (`04-url-fetch.go`) - Fetch and parse remote URLs
5. **Random User-Agent** (`05-random-ua.go`) - Anti-detection with random UAs
6. **Pattern Extraction** (`06-pattern.go`) - Container and field patterns
7. **Chainable Patterns** (`07-chainable-pattern.go`) - Pattern API with fluent syntax

### Running Basic Examples

```bash
# Run a specific example by number
go run examples/basic/*.go -type=01

# Run a specific example by name
go run examples/basic/*.go -type=xpath
go run examples/basic/*.go -type=chainable

# Run all basic examples in sequence
go run examples/basic/*.go -type=all

# Or use Makefile shortcuts
make example-xpath
make example-json
make example-builder
make example-url
make example-random-ua
make example-pattern
make example-chainable
make examples  # Run all
```

## Advanced Examples

The `examples/advanced/` directory contains 9 comprehensive examples demonstrating advanced features:

1. **Basic Product Scraping** (`01-basic-product-scraping.go`) - Container-based extraction with data transformation pipes
2. **XPath Validation** (`02-xpath-validation.go`) - Validate XPath patterns before scraping
3. **Data Cleaning Pipes** (`03-data-cleaning-pipes.go`) - Trim, decode, replace, case transformations
4. **Alternative Patterns** (`04-alternative-patterns.go`) - Fallback patterns for robust extraction
5. **XML Parsing** (`05-xml-parsing.go`) - Sitemap and RSS feed parsing
6. **Real-World E-commerce** (`06-real-world-ecommerce.go`) - Complete product scraping with analytics
7. **URL Health Check** (`07-url-health-check.go`) - Check URL availability concurrently
8. **Configuration Options** (`08-configuration-options.go`) - Suppress errors, timeouts, retries
9. **Custom Pipes** (`09-custom-pipes.go`) - Create and use custom transformation pipes

### Running Advanced Examples

```bash
# Run a specific example by number
go run examples/advanced/*.go -type=01

# Run a specific example by name
go run examples/advanced/*.go -type=basic-product-scraping

# Run all advanced examples in sequence
go run examples/advanced/*.go -type=all

# Show usage and available examples
go run examples/advanced/*.go

# Or use Makefile shortcuts
make example-advanced TYPE=01                    # Run specific example
make example-advanced TYPE=basic-product-scraping  # By name
make examples-advanced                             # Run all
```

## Make Commands

### Build & Test
| Command | Description |
|---------|-------------|
| `make build` | Build library (with vet checks + warnings) |
| `make build-quick` | Quick build (silent, no checks) |
| `make build-verbose` | Build with full verbose output |
| `make test` | Run tests |
| `make test-coverage` | Run tests with coverage report |

### Basic Examples
| Command | Description |
|---------|-------------|
| `make example-xpath` | Run XPath example |
| `make example-json` | Run JSON conversion example |
| `make example-builder` | Run builder/fluent API example |
| `make example-url` | Run URL fetch example |
| `make example-random-ua` | Run Random User-Agent example |
| `make example-pattern` | Run pattern extraction example |
| `make example-chainable` | Run chainable pattern API example |
| `make examples` | Run all basic examples |

### Advanced Examples
| Command | Description |
|---------|-------------|
| `make example-advanced TYPE=<type>` | Run specific advanced example (01-09 or name) |
| `make examples-advanced` | Run all advanced examples |

### Maintenance
| Command | Description |
|---------|-------------|
| `make clean` | Clean build artifacts |
| `make fmt` | Format code |
| `make lint` | Run linter |
| `make deps` | Install dependencies |
| `make release` | Prepare release build |
| `make help` | Show all available commands |
