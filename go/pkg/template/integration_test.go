package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"nof0-api/pkg/llm"
	"nof0-api/pkg/template"
)

// TestJetIntegration tests the full integration of Jet templates via pkg/llm.
func TestJetIntegration(t *testing.T) {
	// Create a temporary directory for templates
	tmpDir := t.TempDir()

	// Create a simple test template
	tmplContent := `Hello, {{.Name}}!
Your balance is {{formatCurrency(.Balance)}}.
{{if isBullish(.Price, .EMA)}}
Market is bullish 🟢
{{else}}
Market is bearish 🔴
{{end}}`

	tmplPath := filepath.Join(tmpDir, "test.jet")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create a prompt template using pkg/llm
	promptTmpl, err := llm.NewPromptTemplate(tmplPath, nil)
	if err != nil {
		t.Fatalf("Failed to create prompt template: %v", err)
	}

	// Test data
	data := map[string]interface{}{
		"Name":    "Alice",
		"Balance": 12500.50,
		"Price":   100.0,
		"EMA":     95.0,
	}

	// Render
	result, err := promptTmpl.Render(data)
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

	// Test digest
	digest := promptTmpl.Digest()
	if len(digest) != 64 { // SHA256 hex string length
		t.Errorf("Expected digest length 64, got %d", len(digest))
	}

	// Test reload
	if err := promptTmpl.Reload(); err != nil {
		t.Errorf("Failed to reload template: %v", err)
	}

	// Re-render after reload
	result2, err := promptTmpl.Render(data)
	if err != nil {
		t.Fatalf("Failed to render after reload: %v", err)
	}

	if result != result2 {
		t.Error("Render results differ after reload")
	}
}

// TestCustomFunctions tests custom functions work through pkg/llm.
func TestCustomFunctions(t *testing.T) {
	tmpDir := t.TempDir()

	// Test each custom function
	tests := []struct {
		name     string
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "formatCurrency",
			template: "{{formatCurrency(.Value)}}",
			data:     map[string]interface{}{"Value": 1500000.0},
			expected: "$1.50M",
		},
		{
			name:     "formatPercent",
			template: "{{formatPercent(.Value)}}",
			data:     map[string]interface{}{"Value": 5.25},
			expected: "+5.25%",
		},
		{
			name:     "colorCode",
			template: "{{colorCode(.Sentiment)}}",
			data:     map[string]interface{}{"Sentiment": "bullish"},
			expected: "🟢",
		},
		{
			name:     "isBullish",
			template: "{{if isBullish(.Price, .EMA)}}yes{{else}}no{{end}}",
			data:     map[string]interface{}{"Price": 100.0, "EMA": 95.0},
			expected: "yes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmplPath := filepath.Join(tmpDir, tt.name+".jet")
			if err := os.WriteFile(tmplPath, []byte(tt.template), 0644); err != nil {
				t.Fatalf("Failed to write template: %v", err)
			}

			promptTmpl, err := llm.NewPromptTemplate(tmplPath, nil)
			if err != nil {
				t.Fatalf("Failed to create prompt template: %v", err)
			}

			result, err := promptTmpl.Render(tt.data)
			if err != nil {
				t.Fatalf("Failed to render template: %v", err)
			}

			if !contains(result, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestBackwardCompatibility ensures the API remains compatible.
func TestBackwardCompatibility(t *testing.T) {
	tmpDir := t.TempDir()

	tmplPath := filepath.Join(tmpDir, "compat.jet")
	if err := os.WriteFile(tmplPath, []byte("Hello, {{.Name}}!"), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Test with nil funcs (backward compatibility)
	promptTmpl1, err := llm.NewPromptTemplate(tmplPath, nil)
	if err != nil {
		t.Fatalf("Failed to create template with nil funcs: %v", err)
	}

	result1, err := promptTmpl1.Render(map[string]interface{}{"Name": "World"})
	if err != nil {
		t.Fatalf("Failed to render: %v", err)
	}

	if !contains(result1, "Hello, World!") {
		t.Errorf("Unexpected output: %q", result1)
	}

	// Test with custom funcs
	customFuncs := map[string]interface{}{
		"upper": func(s string) string {
			return string([]byte(s)) // simplified for test
		},
	}

	promptTmpl2, err := llm.NewPromptTemplate(tmplPath, customFuncs)
	if err != nil {
		t.Fatalf("Failed to create template with custom funcs: %v", err)
	}

	result2, err := promptTmpl2.Render(map[string]interface{}{"Name": "World"})
	if err != nil {
		t.Fatalf("Failed to render with custom funcs: %v", err)
	}

	if !contains(result2, "Hello, World!") {
		t.Errorf("Unexpected output with custom funcs: %q", result2)
	}
}

// TestDocGeneration tests the documentation generator.
func TestDocGeneration(t *testing.T) {
	type TestStruct struct {
		Name    string  `json:"name" doc:"User name" example:"Alice"`
		Balance float64 `json:"balance" doc:"Account balance" example:"1000.50"`
		Active  bool    `json:"active" doc:"Is active" example:"true"`
	}

	gen := template.NewDocGenerator()

	doc, err := gen.Generate(&TestStruct{})
	if err != nil {
		t.Fatalf("Failed to generate documentation: %v", err)
	}

	if doc.Name != "TestStruct" {
		t.Errorf("Expected type name TestStruct, got %s", doc.Name)
	}

	if len(doc.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(doc.Fields))
	}

	// Verify field documentation
	nameField := doc.Fields[0]
	if nameField.Name != "Name" {
		t.Errorf("Expected field Name, got %s", nameField.Name)
	}
	if nameField.JSONName != "name" {
		t.Errorf("Expected JSON name 'name', got %s", nameField.JSONName)
	}
	if nameField.Description != "User name" {
		t.Errorf("Expected description 'User name', got %s", nameField.Description)
	}

	// Test Markdown export
	markdown, err := gen.ExportMarkdown(doc)
	if err != nil {
		t.Fatalf("Failed to export markdown: %v", err)
	}

	if len(markdown) == 0 {
		t.Error("Markdown output is empty")
	}

	t.Logf("Generated markdown:\n%s", markdown)

	// Verify markdown contains expected elements
	expectedParts := []string{
		"TestStruct",
		"Name",
		"Balance",
		"Active",
		"name",
		"balance",
		"active",
	}

	for _, part := range expectedParts {
		if !contains(markdown, part) {
			t.Errorf("Markdown missing expected part: %q", part)
		}
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
