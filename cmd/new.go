package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"api-scaffold/internal/generator"
	"api-scaffold/internal/templating"
	"api-scaffold/templates"
)

var (
	newModule    string
	newAuth      string
	newNoSwagger bool
	newForce     bool
)

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Scaffold a new Gin REST API project",
	Long: `Generates a ready-to-run Gin REST API project with routes, models,
middleware, and (by default) JWT authentication, an OpenAPI spec, and a
Redoc docs viewer.`,
	Example: `  api-scaffold new my-api
  api-scaffold new my-api --module github.com/acme/my-api --auth none
  api-scaffold new my-api --no-swagger`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		if newAuth != "jwt" && newAuth != "none" {
			return fmt.Errorf("--auth must be \"jwt\" or \"none\", got %q", newAuth)
		}

		module := newModule
		if module == "" {
			module = "example/" + projectName
		}

		engine := templating.New(templates.FS)
		opts := generator.ProjectOptions{
			ProjectName: projectName,
			OutputDir:   projectName,
			ModulePath:  module,
			AuthMode:    newAuth,
			SwaggerOn:   !newNoSwagger,
			Force:       newForce,
		}

		if err := generator.NewProject(engine, opts, time.Now().Year()); err != nil {
			return err
		}

		fmt.Printf("Scaffolded %s in ./%s\n\n", projectName, projectName)
		fmt.Println("Next steps:")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Println("  go mod tidy")
		fmt.Println("  go run .")
		if opts.SwaggerOn {
			fmt.Println("\nOnce running, view the API docs at http://localhost:8080/docs/redoc.html")
		}

		return nil
	},
}

func init() {
	newCmd.Flags().StringVar(&newModule, "module", "", "Go module path for the generated project (default \"example/<project-name>\")")
	newCmd.Flags().StringVar(&newAuth, "auth", "jwt", `Authentication mode: "jwt" or "none"`)
	newCmd.Flags().BoolVar(&newNoSwagger, "no-swagger", false, "Skip generating the OpenAPI spec and Redoc docs viewer")
	newCmd.Flags().BoolVar(&newForce, "force", false, "Scaffold into an existing, non-empty directory")

	rootCmd.AddCommand(newCmd)
}
