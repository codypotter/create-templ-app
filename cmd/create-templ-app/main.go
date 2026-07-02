// Command create-templ-app scaffolds a new Go + templ + Tailwind + esbuild web
// app by copying github.com/codypotter/create-templ-app itself and rewriting
// its module path throughout.
//
// Usage:
//
//	go run github.com/codypotter/create-templ-app/cmd/create-templ-app@latest <module-path> [dir]
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const srcModule = "github.com/codypotter/create-templ-app"

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 || strings.HasPrefix(os.Args[1], "-") {
		usage()
		os.Exit(2)
	}

	dstModule := strings.TrimRight(os.Args[1], "/")
	if dstModule == "" {
		fmt.Fprintln(os.Stderr, "create-templ-app: module path must not be empty")
		os.Exit(2)
	}

	dir := path.Base(dstModule)
	if len(os.Args) == 3 {
		dir = os.Args[2]
	}

	if err := run(dstModule, dir); err != nil {
		fmt.Fprintf(os.Stderr, "create-templ-app: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(`Created %s in ./%s

Next steps:
  cd %s
  make tidy    # go mod tidy && npm install

See %s/README.md for templ generate / npm run build / air, etc.
`, dstModule, dir, dir, dir)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: create-templ-app <module-path> [dir]")
}

func run(dstModule, dir string) error {
	if entries, err := os.ReadDir(dir); err == nil {
		if len(entries) > 0 {
			return fmt.Errorf("%s already exists and is not empty", dir)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check target dir: %w", err)
	}

	srcDir, err := resolveSrcDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create target dir: %w", err)
	}

	return copyModuleTree(srcDir, dir, srcModule, dstModule)
}

// resolveSrcDir locates the template's source tree, normally by fetching
// srcModule@latest. Set CREATE_TEMPL_APP_SRC_DIR to copy from a local
// checkout instead, for developing this tool without publishing every
// change first.
func resolveSrcDir() (string, error) {
	if dir := os.Getenv("CREATE_TEMPL_APP_SRC_DIR"); dir != "" {
		return dir, nil
	}

	cmd := exec.Command("go", "mod", "download", "-json", srcModule+"@latest")
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("go mod download: %s", exitErr.Stderr)
		}
		return "", fmt.Errorf("go mod download: %w", err)
	}

	var info struct {
		Dir string
	}
	if err := json.Unmarshal(out, &info); err != nil {
		return "", fmt.Errorf("parse go mod download output: %w", err)
	}
	if info.Dir == "" {
		return "", fmt.Errorf("go mod download: no Dir in output")
	}
	return info.Dir, nil
}

// copyModuleTree copies every regular file under src into dst, replacing
// oldModule with newModule in each file's contents. The replace runs
// uniformly across all files rather than special-casing by extension: it
// rewrites go.mod's module line and every .go import, and is a harmless
// no-op on files with no match (go.sum, static assets, etc.).
func copyModuleTree(src, dst, oldModule, newModule string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			if rel == "." {
				return nil
			}
			return os.MkdirAll(target, 0o755)
		}

		data, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("read %s: %w", rel, err)
		}
		data = bytes.ReplaceAll(data, []byte(oldModule), []byte(newModule))
		if err := os.WriteFile(target, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", rel, err)
		}
		return nil
	})
}
