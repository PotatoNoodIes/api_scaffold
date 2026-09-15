package generator

import (
	"fmt"

	"api-scaffold/internal/config"
	"api-scaffold/internal/templating"
)

const (
	routesSentinel  = "// api-scaffold:routes (do not remove this comment)"
	openapiSentinel = "# api-scaffold:paths (do not remove this comment)"
)

// AddResource generates a model, handler, and routes for a new resource
// inside the scaffolded project found by walking up from startDir, and
// patches its router (and, if enabled, its OpenAPI spec) to wire them in.
func AddResource(engine *templating.Engine, startDir, name string) error {
	root, err := config.FindRoot(startDir)
	if err != nil {
		return err
	}

	cfg, err := config.Load(root)
	if err != nil {
		return err
	}

	if cfg.HasResource(name) {
		return fmt.Errorf("resource %q has already been added to this project", name)
	}

	nameCamel := templating.ToCamel(name)
	nameLower := templating.ToLowerFirst(nameCamel)

	data := templating.ResourceData{
		ProjectData: templating.ProjectData{
			ModulePath: cfg.Module,
			AuthMode:   cfg.AuthMode,
			SwaggerOn:  cfg.SwaggerOn,
		},
		Name:       name,
		NameCamel:  nameCamel,
		NameLower:  nameLower,
		NamePlural: templating.ToPlural(nameLower),
	}

	model, err := engine.Render("resource/model.go.tmpl", data)
	if err != nil {
		return err
	}
	if err := WriteFile(root, fmt.Sprintf("internal/models/%s.go", nameLower), model); err != nil {
		return err
	}

	handler, err := engine.Render("resource/handler.go.tmpl", data)
	if err != nil {
		return err
	}
	if err := WriteFile(root, fmt.Sprintf("internal/handlers/%s.go", nameLower), handler); err != nil {
		return err
	}

	routesSnippet, err := engine.Render("resource/routes_snippet.go.tmpl", data)
	if err != nil {
		return err
	}
	routerPath := root + "/internal/router/router.go"
	if err := InsertGoSnippetBeforeSentinel(routerPath, routesSentinel, routesSnippet); err != nil {
		return err
	}

	if cfg.SwaggerOn {
		openapiSnippet, err := engine.Render("resource/openapi_paths.yaml.tmpl", data)
		if err != nil {
			return err
		}
		openapiPath := root + "/docs/openapi.yaml"
		if err := InsertBeforeSentinel(openapiPath, openapiSentinel, openapiSnippet); err != nil {
			return err
		}
	}

	cfg.Resources = append(cfg.Resources, name)
	if err := config.Save(root, cfg); err != nil {
		return err
	}

	return nil
}
