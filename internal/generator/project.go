package generator

import (
	"api-scaffold/internal/templating"
)

// DefaultGoVersion is the `go` directive written into generated go.mod
// files. It is a fixed, broadly-compatible minimum rather than the host's
// exact toolchain version.
const DefaultGoVersion = "1.22"

// ProjectOptions captures everything needed to scaffold a new project.
type ProjectOptions struct {
	ProjectName string
	OutputDir   string
	ModulePath  string
	AuthMode    string // "jwt" | "none"
	SwaggerOn   bool
	Force       bool
}

// NewProject renders every file in the project manifest and writes it under
// opts.OutputDir, creating the directory if needed.
func NewProject(engine *templating.Engine, opts ProjectOptions, year int) error {
	if err := CheckTargetDir(opts.OutputDir, opts.Force); err != nil {
		return err
	}

	data := templating.ProjectData{
		ProjectName: opts.ProjectName,
		ModulePath:  opts.ModulePath,
		AuthMode:    opts.AuthMode,
		SwaggerOn:   opts.SwaggerOn,
		GoVersion:   DefaultGoVersion,
		Year:        year,
	}

	for _, entry := range ProjectManifest(data) {
		content, err := render(engine, entry, data)
		if err != nil {
			return err
		}
		if err := WriteFile(opts.OutputDir, entry.Output, content); err != nil {
			return err
		}
	}

	return nil
}

func render(engine *templating.Engine, entry ManifestEntry, data any) ([]byte, error) {
	if entry.Static {
		return engine.ReadStatic(entry.Template)
	}
	return engine.Render(entry.Template, data)
}
