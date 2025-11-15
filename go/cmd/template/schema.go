package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"nof0-api/pkg/template"
)

func newSchemaCmd() *cobra.Command {
	var (
		output string
		format string
	)

	cmd := &cobra.Command{
		Use:   "schema [type-name]",
		Short: "Generate documentation schema from Go struct",
		Long: `Generate documentation for a Go struct type.

The schema includes field names, types, JSON names, descriptions,
examples, and required flags extracted from struct tags.

Example:
  template schema SystemPromptData --output=schema.md
  template schema UserPromptData --format=markdown`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			typeName := args[0]

			// Get the type instance
			typeInstance, err := getTypeByName(typeName)
			if err != nil {
				return fmt.Errorf("failed to get type %q: %w", typeName, err)
			}

			// Generate documentation
			gen := template.NewDocGenerator()
			doc, err := gen.Generate(typeInstance)
			if err != nil {
				return fmt.Errorf("failed to generate schema: %w", err)
			}

			// Export based on format
			var result string
			switch format {
			case "markdown", "md":
				result, err = gen.ExportMarkdown(doc)
				if err != nil {
					return fmt.Errorf("failed to export markdown: %w", err)
				}
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}

			// Output
			if output == "" || output == "-" {
				fmt.Print(result)
			} else {
				if err := os.WriteFile(output, []byte(result), 0644); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
				fmt.Fprintf(os.Stderr, "Schema written to: %s\n", output)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
	cmd.Flags().StringVarP(&format, "format", "f", "markdown", "Output format (markdown)")

	return cmd
}

// getTypeByName returns a type instance by name.
// This is a registry of known types that can be documented.
func getTypeByName(name string) (interface{}, error) {
	// Registry of types
	// Note: In production, this should be dynamically populated
	// or read from package metadata
	types := map[string]interface{}{
		// Add your types here as they are defined
		// Example:
		// "SystemPromptData": &examples.SystemPromptData{},
		// "UserPromptData": &examples.UserPromptData{},
	}

	if t, ok := types[name]; ok {
		return t, nil
	}

	return nil, fmt.Errorf("unknown type: %s\n\nAvailable types:\n%s",
		name, getAvailableTypes(types))
}

func getAvailableTypes(types map[string]interface{}) string {
	result := ""
	for name, v := range types {
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		result += fmt.Sprintf("  - %s (%s)\n", name, t.String())
	}
	return result
}
