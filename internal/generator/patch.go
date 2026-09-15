package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
)

// InsertBeforeSentinel inserts snippet immediately above the line containing
// sentinel inside the file at path, rewriting the file in place.
func InsertBeforeSentinel(path, sentinel string, snippet []byte) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("generator: read %s: %w", path, err)
	}

	idx := bytes.Index(content, []byte(sentinel))
	if idx == -1 {
		return fmt.Errorf("generator: sentinel %q not found in %s", sentinel, path)
	}

	lineStart := bytes.LastIndexByte(content[:idx], '\n') + 1

	var out bytes.Buffer
	out.Write(content[:lineStart])
	out.Write(snippet)
	out.Write(content[lineStart:])

	if err := os.WriteFile(path, out.Bytes(), 0o644); err != nil {
		return fmt.Errorf("generator: write %s: %w", path, err)
	}

	return nil
}

// InsertGoSnippetBeforeSentinel behaves like InsertBeforeSentinel but also
// runs the result through gofmt, so a malformed splice surfaces as an error
// instead of corrupting the file.
func InsertGoSnippetBeforeSentinel(path, sentinel string, snippet []byte) error {
	if err := InsertBeforeSentinel(path, sentinel, snippet); err != nil {
		return err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("generator: read %s: %w", path, err)
	}

	formatted, err := format.Source(content)
	if err != nil {
		return fmt.Errorf("generator: %s is not valid Go after patching: %w", path, err)
	}

	if err := os.WriteFile(path, formatted, 0o644); err != nil {
		return fmt.Errorf("generator: write %s: %w", path, err)
	}

	return nil
}
