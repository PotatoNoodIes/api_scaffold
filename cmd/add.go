package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"api-scaffold/internal/generator"
	"api-scaffold/internal/templating"
	"api-scaffold/templates"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add things to an already-scaffolded project",
}

var addResourceCmd = &cobra.Command{
	Use:   "resource <name>",
	Short: "Add a model, handler, and routes for a new resource",
	Long: `Generates a model, an in-memory-backed handler (list/get/create), and
routes for <name>, then wires them into the project's router (and OpenAPI
spec, if enabled). Must be run from inside a project previously created
with "api-scaffold new".`,
	Example: `  api-scaffold add resource widget`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		engine := templating.New(templates.FS)
		if err := generator.AddResource(engine, cwd, name); err != nil {
			return err
		}

		fmt.Printf("Added resource %q (model, handler, routes)\n", name)
		return nil
	},
}

func init() {
	addCmd.AddCommand(addResourceCmd)
	rootCmd.AddCommand(addCmd)
}
