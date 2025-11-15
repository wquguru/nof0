// Package template provides a production-grade template system based on Jet engine.
//
// This package offers:
//   - Jet template engine integration
//   - Custom function library for common formatting needs
//   - Schema documentation generation from Go structs
//   - Template rendering with type safety
package template

import "github.com/CloudyKit/jet/v6"

// Engine represents a template engine that can load and render templates.
type Engine interface {
	// Load loads a template from the specified path
	Load(path string) (*Template, error)

	// Render renders a template with the given data
	Render(tmpl *Template, data interface{}) (string, error)

	// AddFunc adds a custom function to the engine
	AddFunc(name string, fn interface{})
}

// Template represents a loaded template.
type Template struct {
	Name    string
	Path    string
	Content string
	jet     *jet.Template
}

// FieldDoc represents documentation for a struct field.
type FieldDoc struct {
	Name        string      // Field name in Go
	JSONName    string      // JSON field name
	Type        string      // Go type
	Description string      // Field description
	Example     interface{} // Example value
	Required    bool        // Whether field is required
}

// TypeDoc represents documentation for a struct type.
type TypeDoc struct {
	Name        string
	Description string
	Fields      []FieldDoc
}

// DocGenerator generates documentation from Go structs.
type DocGenerator interface {
	// Generate generates documentation for a given struct
	Generate(v interface{}) (*TypeDoc, error)

	// ExportMarkdown exports documentation as Markdown
	ExportMarkdown(doc *TypeDoc) (string, error)
}
