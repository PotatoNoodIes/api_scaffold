# api-scaffold

A terminal-based CLI, built in Go with [Cobra](https://github.com/spf13/cobra),
that scaffolds ready-to-run Gin REST API projects — routes, models,
middleware, JWT authentication, an OpenAPI spec, and a Redoc docs viewer —
so teams don't hand-write the same boilerplate every time.

Every generated file comes from a template rendered with Go's `text/template`,
embedded into the binary via `embed` — the CLI ships as a single static
executable with no runtime file dependencies.

## Install

```sh
go build -o api-scaffold .
```

(Requires Go 1.22+.)

## Usage

### Scaffold a new project

```sh
api-scaffold new my-api
cd my-api
go mod tidy
go run .
```

This generates a Gin app with:

- `GET /health` — liveness check
- `POST /login` / `GET /me` — JWT auth demo (username `demo`, password `demo123`)
- `docs/openapi.yaml` + `docs/redoc.html` — API docs, served at `/docs/redoc.html`

Flags:

| Flag           | Default            | Description                                  |
|----------------|---------------------|-----------------------------------------------|
| `--module`     | `example/<name>`    | Go module path for the generated project      |
| `--auth`       | `jwt`               | `jwt` or `none`                               |
| `--no-swagger` | `false`             | Skip the OpenAPI spec and Redoc docs viewer   |
| `--force`      | `false`             | Scaffold into an existing, non-empty directory|

### Add a resource to an existing project

From inside a project created by `new`:

```sh
api-scaffold add resource widget
```

Generates a model, an in-memory-backed handler (list/get/create), and routes
for `widget`, and wires them into the router (and OpenAPI spec, if enabled)
by inserting above a marker comment and re-formatting with `gofmt`.

### List templates and project state

```sh
api-scaffold list
```

Shows available templates/auth modes, and — if run inside a scaffolded
project — its module path, auth mode, and resources added so far.

## How it works

- `internal/templating` — a small engine (`text/template` + a casing/pluralizing
  `FuncMap`) that renders any embedded template against a data struct.
- `internal/generator` — orchestrates `new` (render the full project manifest)
  and `add resource` (render new files, then patch `router.go` / `openapi.yaml`
  by inserting above a sentinel comment).
- `templates/` — one `.tmpl` file per generated output file, embedded via
  `go:embed`, plus a pinned `redoc.standalone.js` bundled for fully offline
  docs viewing.
- `.api-scaffold.yaml` — written into every generated project, recording its
  module path, auth mode, and resources, so `add` never needs to re-ask.

## Learn more

`sites/index.html` is a self-contained explainer page for this project —
purpose, quickstart, command reference, generated output, and architecture.
Open it directly in a browser, or serve it locally:

```sh
cd sites && python3 -m http.server 8000
```

## Testing

```sh
go test ./...
```

`test/e2e_test.go` builds the CLI, scaffolds real projects into a temp
directory, builds and runs them, and exercises their HTTP endpoints — the
strongest signal that generated code is actually valid and runnable.

## License

MIT — see [LICENSE](LICENSE).
