package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"nof0-api/pkg/template"
)

func newDocCmd() *cobra.Command {
	var (
		format string
	)

	cmd := &cobra.Command{
		Use:   "doc [type-name]",
		Short: "Display documentation for a type",
		Long: `Display documentation for a Go struct type in terminal format.

This command shows type information in a human-readable format
suitable for terminal display.

Example:
  template doc SystemPromptData
  template doc UserPromptData --format=table`,
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
				return fmt.Errorf("failed to generate documentation: %w", err)
			}

			// Display based on format
			switch format {
			case "table":
				displayTable(doc)
			case "simple":
				displaySimple(doc)
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, simple)")

	return cmd
}

func displayTable(doc *template.TypeDoc) {
	fmt.Printf("Type: %s\n", doc.Name)
	if doc.Description != "" {
		fmt.Printf("Description: %s\n", doc.Description)
	}
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FIELD\tTYPE\tJSON\tDESCRIPTION\tEXAMPLE")
	fmt.Fprintln(w, strings.Repeat("-", 80))

	for _, field := range doc.Fields {
		required := ""
		if field.Required {
			required = " *"
		}

		jsonName := field.JSONName
		if jsonName == "" {
			jsonName = "-"
		}

		example := formatExampleForDisplay(field.Example)

		fmt.Fprintf(w, "%s%s\t%s\t%s\t%s\t%s\n",
			field.Name,
			required,
			field.Type,
			jsonName,
			field.Description,
			example,
		)
	}

	w.Flush()
	fmt.Println()
	fmt.Println("* = required field")
}

func displaySimple(doc *template.TypeDoc) {
	fmt.Printf("Type: %s\n", doc.Name)
	if doc.Description != "" {
		fmt.Printf("Description: %s\n", doc.Description)
	}
	fmt.Println()

	for _, field := range doc.Fields {
		required := ""
		if field.Required {
			required = " (required)"
		}

		fmt.Printf("%s%s\n", field.Name, required)
		fmt.Printf("  Type: %s\n", field.Type)
		if field.JSONName != "" && field.JSONName != "-" {
			fmt.Printf("  JSON: %s\n", field.JSONName)
		}
		if field.Description != "" {
			fmt.Printf("  Description: %s\n", field.Description)
		}
		if field.Example != nil && field.Example != "" {
			fmt.Printf("  Example: %v\n", field.Example)
		}
		fmt.Println()
	}
}

func formatExampleForDisplay(example interface{}) string {
	if example == nil || example == "" {
		return "-"
	}

	str := fmt.Sprintf("%v", example)
	if len(str) > 30 {
		return str[:27] + "..."
	}
	return str
}
