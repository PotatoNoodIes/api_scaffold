package generator

import "api-scaffold/internal/templating"

// ManifestEntry maps a single embedded template (or static asset) to an
// output path relative to the generated project's root.
type ManifestEntry struct {
	Template string // path within templates.FS, e.g. "project/main.go.tmpl"
	Output   string // path relative to the project root, e.g. "main.go"
	Static   bool   // if true, copy the template's bytes verbatim instead of rendering it
}

// ProjectManifest returns the set of files to generate for `new`, given the
// project's configuration. Auth and swagger files are included only when
// requested, keeping the output minimal by default.
func ProjectManifest(data templating.ProjectData) []ManifestEntry {
	entries := []ManifestEntry{
		{Template: "project/go.mod.tmpl", Output: "go.mod"},
		{Template: "project/main.go.tmpl", Output: "main.go"},
		{Template: "project/README.md.tmpl", Output: "README.md"},
		{Template: "project/.env.example.tmpl", Output: ".env.example"},
		{Template: "project/.gitignore.tmpl", Output: ".gitignore"},
		{Template: "project/.api-scaffold.yaml.tmpl", Output: ".api-scaffold.yaml"},
		{Template: "project/internal/config/config.go.tmpl", Output: "internal/config/config.go"},
		{Template: "project/internal/router/router.go.tmpl", Output: "internal/router/router.go"},
		{Template: "project/internal/handlers/health.go.tmpl", Output: "internal/handlers/health.go"},
		{Template: "project/internal/middleware/logging.go.tmpl", Output: "internal/middleware/logging.go"},
	}

	if data.AuthMode == "jwt" {
		entries = append(entries,
			ManifestEntry{Template: "project/internal/models/user.go.tmpl", Output: "internal/models/user.go"},
			ManifestEntry{Template: "project/internal/middleware/auth.go.tmpl", Output: "internal/middleware/auth.go"},
			ManifestEntry{Template: "project/internal/handlers/auth.go.tmpl", Output: "internal/handlers/auth.go"},
		)
	}

	if data.SwaggerOn {
		entries = append(entries,
			ManifestEntry{Template: "project/docs/openapi.yaml.tmpl", Output: "docs/openapi.yaml"},
			ManifestEntry{Template: "project/docs/redoc.html.tmpl", Output: "docs/redoc.html"},
			ManifestEntry{Template: "static/redoc.standalone.js", Output: "docs/redoc.standalone.js", Static: true},
		)
	}

	return entries
}
