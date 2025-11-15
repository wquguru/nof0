package template

import (
	"fmt"
	"reflect"
	"strings"
)

// SimpleDocGenerator generates documentation from Go structs.
type SimpleDocGenerator struct{}

// NewDocGenerator creates a new documentation generator.
func NewDocGenerator() *SimpleDocGenerator {
	return &SimpleDocGenerator{}
}

// Generate generates documentation for a struct.
func (g *SimpleDocGenerator) Generate(v interface{}) (*TypeDoc, error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Handle pointer
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", typ.Kind())
	}

	doc := &TypeDoc{
		Name:        typ.Name(),
		Description: extractTypeDoc(typ),
		Fields:      make([]FieldDoc, 0, typ.NumField()),
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldDoc := FieldDoc{
			Name:        field.Name,
			JSONName:    extractJSONName(field),
			Type:        field.Type.String(),
			Description: extractFieldDoc(field),
			Example:     extractExample(field),
			Required:    isRequired(field),
		}

		doc.Fields = append(doc.Fields, fieldDoc)
	}

	return doc, nil
}

// ExportMarkdown exports documentation as Markdown.
func (g *SimpleDocGenerator) ExportMarkdown(doc *TypeDoc) (string, error) {
	var buf strings.Builder

	// Title
	buf.WriteString(fmt.Sprintf("# %s\n\n", doc.Name))

	if doc.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", doc.Description))
	}

	// Table header
	buf.WriteString("| Field | Type | Template Variable | Description | Example |\n")
	buf.WriteString("|-------|------|-------------------|-------------|----------|\n")

	// Table rows
	for _, field := range doc.Fields {
		templateVar := fmt.Sprintf("{{.%s}}", field.Name)
		if field.JSONName != "" && field.JSONName != "-" {
			templateVar = fmt.Sprintf("{{.%s}} or {{.%s}}", field.Name, field.JSONName)
		}

		required := ""
		if field.Required {
			required = "✓ "
		}

		example := formatExample(field.Example)

		buf.WriteString(fmt.Sprintf("| %s | %s | `%s` | %s%s | `%s` |\n",
			field.Name,
			field.Type,
			templateVar,
			required,
			field.Description,
			example,
		))
	}

	return buf.String(), nil
}

// extractJSONName extracts JSON field name from struct tag.
func extractJSONName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return ""
	}

	// Handle "name,omitempty"
	parts := strings.Split(tag, ",")
	return parts[0]
}

// extractFieldDoc extracts field documentation from struct tag or comment.
func extractFieldDoc(field reflect.StructField) string {
	// Try "doc" tag first
	if doc := field.Tag.Get("doc"); doc != "" {
		return doc
	}

	// Try "description" tag
	if desc := field.Tag.Get("description"); desc != "" {
		return desc
	}

	return ""
}

// extractExample extracts example value from struct tag.
func extractExample(field reflect.StructField) interface{} {
	if example := field.Tag.Get("example"); example != "" {
		return example
	}
	return ""
}

// isRequired checks if field is marked as required.
func isRequired(field reflect.StructField) bool {
	schema := field.Tag.Get("schema")
	return strings.Contains(schema, "required")
}

// extractTypeDoc extracts type-level documentation.
func extractTypeDoc(typ reflect.Type) string {
	// This would require parsing Go comments, which is complex
	// For now, return empty string
	// Can be enhanced with go/ast parsing
	return ""
}

// formatExample formats an example value for display.
func formatExample(example interface{}) string {
	if example == nil || example == "" {
		return ""
	}

	switch v := example.(type) {
	case string:
		if v == "" {
			return ""
		}
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
