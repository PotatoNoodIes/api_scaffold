// Package test contains end-to-end tests that build the api-scaffold CLI,
// run it to generate real projects, and verify the generated projects
// actually compile and serve traffic correctly. This is the strongest
// signal that the templating system produces valid, runnable Go code.
package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// buildCLI compiles the api-scaffold binary into dir and returns its path.
func buildCLI(t *testing.T, dir string) string {
	t.Helper()

	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	bin := filepath.Join(dir, "api-scaffold")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = repoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build api-scaffold: %v\n%s", err, out)
	}

	return bin
}

func run(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
	return string(out)
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func waitHealthy(t *testing.T, base string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(base + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("server at %s never became healthy", base)
}

func startServer(t *testing.T, projectDir string, port int) func() {
	t.Helper()

	bin := filepath.Join(projectDir, "app-under-test")
	buildCmd := exec.Command("go", "build", "-o", bin, ".")
	buildCmd.Dir = projectDir
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("build generated project: %v\n%s", err, out)
	}

	cmd := exec.Command(bin)
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start generated project: %v", err)
	}

	base := fmt.Sprintf("http://127.0.0.1:%d", port)
	waitHealthy(t, base)

	return func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}
}

func TestNewJWTProject_BuildsRunsAndAuthenticates(t *testing.T) {
	tmp := t.TempDir()
	cli := buildCLI(t, tmp)

	workDir := filepath.Join(tmp, "work")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatal(err)
	}
	run(t, workDir, cli, "new", "jwtapp", "--module", "example/jwtapp")

	projectDir := filepath.Join(workDir, "jwtapp")
	run(t, projectDir, "go", "mod", "tidy")

	port := freePort(t)
	stop := startServer(t, projectDir, port)
	defer stop()

	base := fmt.Sprintf("http://127.0.0.1:%d", port)

	// /me without a token must be rejected.
	resp, err := http.Get(base + "/me")
	if err != nil {
		t.Fatalf("GET /me: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("GET /me without token: got %d, want 401", resp.StatusCode)
	}

	// Login with the demo credentials.
	loginBody := bytes.NewBufferString(`{"username":"demo","password":"demo123"}`)
	resp, err = http.Post(base+"/login", "application/json", loginBody)
	if err != nil {
		t.Fatalf("POST /login: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST /login: got %d, want 200", resp.StatusCode)
	}
	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if loginResp.Token == "" {
		t.Fatal("login response did not contain a token")
	}

	// /me with the token must succeed.
	req, _ := http.NewRequest(http.MethodGet, base+"/me", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /me with token: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /me with token: got %d, want 200", resp.StatusCode)
	}

	// Docs must be servable.
	resp, err = http.Get(base + "/docs/redoc.html")
	if err != nil {
		t.Fatalf("GET /docs/redoc.html: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /docs/redoc.html: got %d, want 200", resp.StatusCode)
	}

	resp, err = http.Get(base + "/docs/openapi.yaml")
	if err != nil {
		t.Fatalf("GET /docs/openapi.yaml: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /docs/openapi.yaml: got %d, want 200", resp.StatusCode)
	}
}

func TestNewNoAuthProject_BuildsAndServesHealthOnly(t *testing.T) {
	tmp := t.TempDir()
	cli := buildCLI(t, tmp)

	workDir := filepath.Join(tmp, "work")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatal(err)
	}
	run(t, workDir, cli, "new", "noauthapp", "--auth", "none", "--module", "example/noauthapp")

	projectDir := filepath.Join(workDir, "noauthapp")
	run(t, projectDir, "go", "mod", "tidy")

	port := freePort(t)
	stop := startServer(t, projectDir, port)
	defer stop()

	base := fmt.Sprintf("http://127.0.0.1:%d", port)

	resp, err := http.Post(base+"/login", "application/json", nil)
	if err != nil {
		t.Fatalf("POST /login: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("POST /login on auth=none project: got %d, want 404", resp.StatusCode)
	}
}

func TestAddResource_WiresRoutesAndOpenAPI(t *testing.T) {
	tmp := t.TempDir()
	cli := buildCLI(t, tmp)

	workDir := filepath.Join(tmp, "work")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatal(err)
	}
	run(t, workDir, cli, "new", "widgetapp", "--auth", "none", "--module", "example/widgetapp")

	projectDir := filepath.Join(workDir, "widgetapp")
	run(t, projectDir, cli, "add", "resource", "widget")
	run(t, projectDir, "go", "mod", "tidy")
	run(t, projectDir, "go", "vet", "./...")

	port := freePort(t)
	stop := startServer(t, projectDir, port)
	defer stop()

	base := fmt.Sprintf("http://127.0.0.1:%d", port)

	createBody := bytes.NewBufferString(`{"name":"sprocket"}`)
	resp, err := http.Post(base+"/widgets", "application/json", createBody)
	if err != nil {
		t.Fatalf("POST /widgets: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST /widgets: got %d, want 201", resp.StatusCode)
	}

	resp, err = http.Get(base + "/widgets")
	if err != nil {
		t.Fatalf("GET /widgets: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /widgets: got %d, want 200", resp.StatusCode)
	}

	// Adding the same resource twice must fail cleanly.
	cmd := exec.Command(cli, "add", "resource", "widget")
	cmd.Dir = projectDir
	if err := cmd.Run(); err == nil {
		t.Fatal("expected re-adding an existing resource to fail")
	}
}
