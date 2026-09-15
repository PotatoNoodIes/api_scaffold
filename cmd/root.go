package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api-scaffold",
	Short: "Scaffold Go REST API projects with routes, models, middleware, and auth",
	Long: `api-scaffold is a CLI tool that generates ready-to-run Go REST API
projects (built on Gin) complete with routes, models, middleware,
JWT authentication, an OpenAPI spec, and a Redoc docs viewer.

It also lets you add new resources to an already-scaffolded project.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
