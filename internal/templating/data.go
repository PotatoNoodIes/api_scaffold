package templating

// ProjectData describes a project being scaffolded by `api-scaffold new`.
type ProjectData struct {
	ProjectName string // "my-api"
	ModulePath  string // "github.com/acme/my-api"
	AuthMode    string // "jwt" | "none"
	SwaggerOn   bool
	GoVersion   string
	Year        int
}

// ResourceData describes a resource being added by `api-scaffold add resource`.
type ResourceData struct {
	ProjectData
	Name       string // "widget" (as given by the user)
	NamePlural string // "widgets"
	NameCamel  string // "Widget"
	NameLower  string // "widget"
}
