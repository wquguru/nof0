# Template System

A production-grade template system for LLM prompts using the Jet template engine.

## Features

- **Jet Template Engine**: Clear Django/Jinja2-like syntax
- **Custom Functions**: Trading-specific helpers (formatCurrency, colorCode, etc.)
- **Documentation Generator**: Auto-generate docs from Go structs
- **CLI Tool**: Command-line utility for template management
- **Backward Compatible**: Drop-in replacement for text/template

## Quick Start

### Using Templates in Code

```go
import "nof0-api/pkg/llm"

// Create a template
tmpl, err := llm.NewPromptTemplate("path/to/template.jet", nil)
if err != nil {
    log.Fatal(err)
}

// Render with data
data := map[string]interface{}{
    "Name": "Alice",
    "Balance": 12500.50,
}

result, err := tmpl.Render(data)
```

### Using the CLI Tool

```bash
# Build the CLI tool
go build -o bin/template ./cmd/template

# Render a template
./bin/template render prompt.jet --data data.json

# Generate documentation for a struct
./bin/template schema MyStruct --output schema.md

# Display type documentation
./bin/template doc MyStruct

# List available types
./bin/template list
```

## Template Syntax

### Variables
```jet
Hello, {{.Name}}!
Your balance is {{.Balance}}.
```

### Comments
```jet
{* This is a comment *}
```

### Conditionals
```jet
{{if .Price > .EMA}}
Market is bullish 🟢
{{else}}
Market is bearish 🔴
{{end}}
```

### Custom Functions
```jet
Balance: {{formatCurrency(.Balance)}}
Change: {{formatPercent(.Change)}}
Status: {{colorCode(.Sentiment)}}
```

## Available Custom Functions

### Formatting
- `formatCurrency(value)` - Format numbers as currency ($1.50M, $5.50K)
- `formatPercent(value)` - Format percentages with +/- sign
- `formatFloat(value, precision)` - Format floats with precision

### Indicators
- `colorCode(sentiment)` - Return emoji indicators (🟢🔴🟡⚪)
- `trendIndicator(current, previous)` - Return trend arrows (📈📉➡️)

### Helpers
- `isBullish(price, ema)` - Check if price is above EMA
- `isBearish(price, ema)` - Check if price is below EMA
- `isOverbought(rsi)` - Check if RSI > 70
- `isOversold(rsi)` - Check if RSI < 30

## Documentation Generation

Add struct tags to enable documentation generation:

```go
type PromptData struct {
    Name     string  `json:"name" doc:"User name" example:"Alice"`
    Balance  float64 `json:"balance" doc:"Account balance" example:"1000.50"`
    Active   bool    `json:"active" doc:"Is account active" example:"true"`
}
```

Generate documentation:

```bash
./bin/template schema PromptData --output docs/prompt-data.md
```

## Migration from text/template

The pkg/llm API remains unchanged, but templates should be migrated to Jet syntax:

### Before (text/template)
```
{{/* Comment */}}
{{printf "%.2f" .Value}}
{{if gt .Price .EMA}}bullish{{end}}
```

### After (Jet)
```
{* Comment *}
{{formatFloat(.Value, 2)}}
{{if .Price > .EMA}}bullish{{end}}
```

Key differences:
- Comments: `{{/* */}}` → `{* *}`
- Comparisons: `{{if gt .x 5}}` → `{{if .x > 5}}`
- No `printf` needed, use custom functions

## File Structure

```
pkg/template/
  ├── types.go           # Core interfaces and types
  ├── jet.go             # Jet engine wrapper
  ├── funcs.go           # Custom function library
  ├── doc.go             # Documentation generator
  ├── jet_test.go        # Unit tests
  └── integration_test.go # Integration tests

cmd/template/
  ├── main.go            # CLI entry point
  ├── schema.go          # Schema generation command
  ├── doc.go             # Documentation display command
  ├── list.go            # List types command
  └── render.go          # Template rendering command

pkg/llm/
  └── prompt.go          # Updated to use pkg/template
```

## Testing

Run all tests:
```bash
go test -v ./pkg/template/...
```

Test specific functionality:
```bash
go test -v ./pkg/template/... -run TestJetIntegration
go test -v ./pkg/template/... -run TestCustomFunctions
go test -v ./pkg/template/... -run TestDocGeneration
```

## Examples

See the integration tests in `pkg/template/integration_test.go` for comprehensive examples of:
- Basic template rendering
- Custom functions
- Backward compatibility
- Documentation generation

## Performance

Jet is approximately 5x faster than text/template and supports:
- Template caching
- Development mode with auto-reload
- Efficient memory usage

## Dependencies

- `github.com/CloudyKit/jet/v6` - Jet template engine
- `github.com/spf13/cobra` - CLI framework (cmd/template only)

## Notes

- Template file extension: `.jet` (recommended) or `.tmpl` (backward compatible)
- Comments use `{* *}` syntax, not `{{/* */}}`
- All variables are accessed via dot notation: `{{.Field}}`
- Custom functions are automatically registered for all templates
- The system maintains backward compatibility with existing pkg/llm usage
