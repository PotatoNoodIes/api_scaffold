package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"api-scaffold/internal/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available scaffold templates, auth modes, and project resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Project templates:")
		fmt.Println("  gin-rest-api   Gin-based REST API with routes, models, middleware, auth")
		fmt.Println()
		fmt.Println("Auth modes (--auth):")
		fmt.Println("  jwt    JWT-based login + protected routes (default)")
		fmt.Println("  none   No authentication scaffolding")

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		root, err := config.FindRoot(cwd)
		if err != nil {
			return nil // not inside a scaffolded project; nothing more to show
		}

		cfg, err := config.Load(root)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Printf("Current project (%s):\n", root)
		fmt.Printf("  module:    %s\n", cfg.Module)
		fmt.Printf("  auth mode: %s\n", cfg.AuthMode)
		fmt.Printf("  swagger:   %v\n", cfg.SwaggerOn)
		if len(cfg.Resources) == 0 {
			fmt.Println("  resources: (none yet -- use `api-scaffold add resource <name>`)")
		} else {
			fmt.Printf("  resources: %v\n", cfg.Resources)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
