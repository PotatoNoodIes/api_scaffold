// Package config reads and writes the `.api-scaffold.yaml` file that every
// scaffolded project carries at its root, recording enough state (module
// path, auth mode, resources already added) that `api-scaffold add` can
// operate without the user re-specifying flags.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileName is the name of the config file written at a scaffolded
// project's root.
const FileName = ".api-scaffold.yaml"

// ScaffoldConfig is the persisted state of a scaffolded project.
type ScaffoldConfig struct {
	Module    string   `yaml:"module"`
	AuthMode  string   `yaml:"auth_mode"`
	SwaggerOn bool     `yaml:"swagger_enabled"`
	Resources []string `yaml:"resources"`
}

// FindRoot walks upward from startDir looking for FileName, returning the
// directory that contains it. It errors if it reaches the filesystem root
// without finding one.
func FindRoot(startDir string) (string, error) {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, FileName)); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not inside an api-scaffold project (no %s found in %s or any parent directory)", FileName, startDir)
		}
		dir = parent
	}
}

// Load reads and parses the config file at dir/FileName.
func Load(dir string) (*ScaffoldConfig, error) {
	data, err := os.ReadFile(filepath.Join(dir, FileName))
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", FileName, err)
	}

	var cfg ScaffoldConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %s: %w", FileName, err)
	}

	return &cfg, nil
}

// Save writes cfg back to dir/FileName.
func Save(dir string, cfg *ScaffoldConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("config: marshal %s: %w", FileName, err)
	}

	if err := os.WriteFile(filepath.Join(dir, FileName), data, 0o644); err != nil {
		return fmt.Errorf("config: write %s: %w", FileName, err)
	}

	return nil
}

// HasResource reports whether name is already recorded in cfg.Resources.
func (cfg *ScaffoldConfig) HasResource(name string) bool {
	for _, r := range cfg.Resources {
		if r == name {
			return true
		}
	}
	return false
}
