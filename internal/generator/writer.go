package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFile writes content to filepath.Join(root, relPath), creating any
// missing parent directories.
func WriteFile(root, relPath string, content []byte) error {
	full := filepath.Join(root, relPath)

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return fmt.Errorf("generator: create dir for %s: %w", relPath, err)
	}

	if err := os.WriteFile(full, content, 0o644); err != nil {
		return fmt.Errorf("generator: write %s: %w", relPath, err)
	}

	return nil
}

// CheckTargetDir ensures dir is safe to scaffold into: it must not exist yet,
// or must exist and be empty, unless force is set.
func CheckTargetDir(dir string, force bool) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0o755)
	}
	if err != nil {
		return fmt.Errorf("generator: stat %s: %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s already exists and is not a directory", dir)
	}
	if force {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("generator: read %s: %w", dir, err)
	}
	if len(entries) > 0 {
		return fmt.Errorf("directory %s already exists and is not empty (use --force to scaffold into it anyway)", dir)
	}

	return nil
}
