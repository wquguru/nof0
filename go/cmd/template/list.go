package main

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available types for documentation",
		Long: `List all registered types that can be used with the schema and doc commands.

Example:
  template list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get all registered types
			types := getRegisteredTypes()

			if len(types) == 0 {
				fmt.Println("No types registered.")
				fmt.Println()
				fmt.Println("To register types, add them to the type registry in schema.go")
				return nil
			}

			fmt.Println("Available types:")
			fmt.Println()

			// Sort type names for consistent output
			names := make([]string, 0, len(types))
			for name := range types {
				names = append(names, name)
			}
			sort.Strings(names)

			for _, name := range names {
				v := types[name]
				t := reflect.TypeOf(v)
				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}

				// Count fields
				numFields := 0
				if t.Kind() == reflect.Struct {
					for i := 0; i < t.NumField(); i++ {
						if t.Field(i).IsExported() {
							numFields++
						}
					}
				}

				fmt.Printf("  %s\n", name)
				fmt.Printf("    Type: %s\n", t.String())
				fmt.Printf("    Fields: %d\n", numFields)
				fmt.Println()
			}

			fmt.Println("Usage:")
			fmt.Println("  template schema <type-name> -o output.md")
			fmt.Println("  template doc <type-name>")

			return nil
		},
	}

	return cmd
}

func getRegisteredTypes() map[string]interface{} {
	// This should return the same registry as in schema.go
	// In a real implementation, you'd have a shared registry
	return map[string]interface{}{
		// Add your types here
		// Example:
		// "SystemPromptData": &examples.SystemPromptData{},
		// "UserPromptData": &examples.UserPromptData{},
	}
}
