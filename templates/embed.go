// Package templates embeds every scaffold template shipped with the
// api-scaffold binary so it can generate projects without needing any
// files on disk at runtime.
package templates

import "embed"

//go:embed all:project all:resource all:static
var FS embed.FS
