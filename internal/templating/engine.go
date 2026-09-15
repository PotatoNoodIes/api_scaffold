package templating

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"text/template"
)

// Engine renders named templates from an underlying filesystem (typically
// an embed.FS baked into the CLI binary) against arbitrary template data.
type Engine struct {
	fsys fs.FS
}

// New builds an Engine backed by fsys.
func New(fsys fs.FS) *Engine {
	return &Engine{fsys: fsys}
}

// Render reads the template at name, executes it with data, and returns the
// rendered bytes. name is a slash-separated path relative to the Engine's
// filesystem root, e.g. "project/main.go.tmpl".
func (e *Engine) Render(name string, data any) ([]byte, error) {
	content, err := fs.ReadFile(e.fsys, name)
	if err != nil {
		return nil, fmt.Errorf("templating: read %s: %w", name, err)
	}

	tmpl, err := template.New(path.Base(name)).Funcs(FuncMap()).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("templating: parse %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("templating: execute %s: %w", name, err)
	}

	return buf.Bytes(), nil
}

// ReadStatic returns the raw, unrendered bytes of a non-template asset
// (e.g. the bundled Redoc JS bundle).
func (e *Engine) ReadStatic(name string) ([]byte, error) {
	content, err := fs.ReadFile(e.fsys, name)
	if err != nil {
		return nil, fmt.Errorf("templating: read static %s: %w", name, err)
	}
	return content, nil
}
