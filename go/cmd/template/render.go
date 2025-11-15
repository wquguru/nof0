package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"nof0-api/pkg/template"
)

func newRenderCmd() *cobra.Command {
	var (
		templateDir string
		dataFile    string
		devMode     bool
	)

	cmd := &cobra.Command{
		Use:   "render [template-file]",
		Short: "Render a Jet template with data",
		Long: `Render a Jet template file with data from a JSON file.

This is useful for testing templates before integrating them
into your application.

Example:
  template render prompt.jet --data data.json
  template render system.jet --data data.json --template-dir ./templates`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateFile := args[0]

			// Read data file
			var data map[string]interface{}
			if dataFile != "" {
				content, err := os.ReadFile(dataFile)
				if err != nil {
					return fmt.Errorf("failed to read data file: %w", err)
				}

				if err := json.Unmarshal(content, &data); err != nil {
					return fmt.Errorf("failed to parse data file: %w", err)
				}
			}

			// Create engine
			engine := template.NewJetEngine(template.JetOptions{
				TemplateDir:     templateDir,
				DevelopmentMode: devMode,
			})

			// Load template
			tmpl, err := engine.Load(templateFile)
			if err != nil {
				return fmt.Errorf("failed to load template: %w", err)
			}

			// Render
			result, err := engine.Render(tmpl, data)
			if err != nil {
				return fmt.Errorf("failed to render template: %w", err)
			}

			fmt.Print(result)

			return nil
		},
	}

	cmd.Flags().StringVar(&templateDir, "template-dir", "./templates", "Template directory")
	cmd.Flags().StringVar(&dataFile, "data", "", "JSON data file")
	cmd.Flags().BoolVar(&devMode, "dev", false, "Enable development mode (auto-reload)")

	return cmd
}
