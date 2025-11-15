package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Template CLI tool for managing Jet templates and documentation",
		Long: `A command-line tool for working with Jet templates.

Features:
  - Generate documentation from Go structs
  - Display type information
  - List available template types
  - Render templates with data`,
		Version: version,
	}

	// Add subcommands
	cmd.AddCommand(newSchemaCmd())
	cmd.AddCommand(newDocCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newRenderCmd())

	return cmd
}
