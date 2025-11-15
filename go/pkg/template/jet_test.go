package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJetEngine(t *testing.T) {
	// Create temp directory for templates
	tmpDir := t.TempDir()

	// Create a test template
	tmplContent := `Hello, {{.Name}}!
Your balance is {{.Balance | formatCurrency}}.
{{if isBullish(.Price, .EMA)}}
Market is bullish 🟢
{{else}}
Market is bearish 🔴
{{end}}`

	tmplPath := filepath.Join(tmpDir, "test.jet")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create engine
	engine := NewJetEngine(JetOptions{
		TemplateDir:     tmpDir,
		DevelopmentMode: true,
	})

	// Load template
	tmpl, err := engine.Load("test.jet")
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// Test data
	data := map[string]interface{}{
		"Name":    "Alice",
		"Balance": 12500.50,
		"Price":   100.0,
		"EMA":     95.0,
	}

	// Render
	result, err := engine.Render(tmpl, data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Verify output
	if len(result) == 0 {
		t.Error("Rendered output is empty")
	}

	t.Logf("Rendered output:\n%s", result)

	// Check for expected content
	expectedParts := []string{
		"Alice",
		"$12.50K", // formatted currency
		"bullish", // Price > EMA
	}

	for _, part := range expectedParts {
		if !contains(result, part) {
			t.Errorf("Output missing expected part: %q", part)
		}
	}
}

func TestCustomFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() string
		expected string
	}{
		{
			name:     "FormatCurrency_Large",
			fn:       func() string { return FormatCurrency(1500000) },
			expected: "$1.50M",
		},
		{
			name:     "FormatCurrency_Medium",
			fn:       func() string { return FormatCurrency(5500) },
			expected: "$5.50K",
		},
		{
			name:     "FormatCurrency_Small",
			fn:       func() string { return FormatCurrency(99.99) },
			expected: "$99.99",
		},
		{
			name:     "FormatPercent_Positive",
			fn:       func() string { return FormatPercent(5.25) },
			expected: "+5.25%",
		},
		{
			name:     "FormatPercent_Negative",
			fn:       func() string { return FormatPercent(-2.50) },
			expected: "-2.50%",
		},
		{
			name:     "ColorCode_Bullish",
			fn:       func() string { return ColorCode("bullish") },
			expected: "🟢",
		},
		{
			name:     "ColorCode_Bearish",
			fn:       func() string { return ColorCode("bearish") },
			expected: "🔴",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() bool
		expected bool
	}{
		{
			name:     "IsBullish_True",
			fn:       func() bool { return IsBullish(100, 95) },
			expected: true,
		},
		{
			name:     "IsBullish_False",
			fn:       func() bool { return IsBullish(90, 95) },
			expected: false,
		},
		{
			name:     "IsOverbought_True",
			fn:       func() bool { return IsOverbought(75) },
			expected: true,
		},
		{
			name:     "IsOverbought_False",
			fn:       func() bool { return IsOverbought(65) },
			expected: false,
		},
		{
			name:     "IsOversold_True",
			fn:       func() bool { return IsOversold(25) },
			expected: true,
		},
		{
			name:     "IsOversold_False",
			fn:       func() bool { return IsOversold(35) },
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && stringContains(s, substr)
}

// stringContains is a helper to check substring
func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
