package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyModuleTree(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "out")

	writeFile(t, filepath.Join(src, "go.mod"), "module github.com/old/app\n\ngo 1.25\n")
	writeFile(t, filepath.Join(src, "go.sum"), "some.dep v1.0.0 h1:abc=\n")
	writeFile(t, filepath.Join(src, "cmd", "server", "main.go"),
		"package main\n\nimport \"github.com/old/app/internal/config\"\n\nfunc main() {}\n")

	if err := copyModuleTree(src, dst, "github.com/old/app", "github.com/new/app"); err != nil {
		t.Fatalf("copyModuleTree: %v", err)
	}

	goMod := readFile(t, filepath.Join(dst, "go.mod"))
	if !strings.Contains(goMod, "module github.com/new/app") {
		t.Errorf("go.mod = %q, want it to contain the rewritten module path", goMod)
	}

	mainGo := readFile(t, filepath.Join(dst, "cmd", "server", "main.go"))
	if !strings.Contains(mainGo, `"github.com/new/app/internal/config"`) {
		t.Errorf("main.go = %q, want it to contain the rewritten import", mainGo)
	}
	if strings.Contains(mainGo, "github.com/old/app") {
		t.Errorf("main.go = %q, still contains old module path", mainGo)
	}

	goSum := readFile(t, filepath.Join(dst, "go.sum"))
	if !strings.Contains(goSum, "some.dep v1.0.0") {
		t.Errorf("go.sum = %q, want unchanged content (no match to rewrite)", goSum)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
