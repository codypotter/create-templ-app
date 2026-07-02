package main

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/mod/module"
	modzip "golang.org/x/mod/zip"
)

// TestSrcModuleExcludesCLI verifies the invariant this whole design leans
// on: because this package has its own go.mod, Go's module-zip packaging
// (the same logic `go mod download` uses to fetch srcModule) omits this
// directory from the zip, which is what keeps the CLI's own code out of
// scaffolded projects.
//
// This runs entirely offline against the actual git-tracked HEAD content —
// no network fetch or published tag required — rather than the live
// working tree, which may have untracked build artifacts (node_modules/,
// dist/) that a real fetch would never see.
func TestSrcModuleExcludesCLI(t *testing.T) {
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}

	trackedDir := t.TempDir()
	extractGitHEAD(t, repoRoot, trackedDir)

	zipPath := filepath.Join(t.TempDir(), "module.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}

	mv := module.Version{Path: srcModule, Version: "v0.0.0-20240101000000-000000000000"}
	err = modzip.CreateFromDir(f, mv, trackedDir)
	f.Close()
	if err != nil {
		t.Fatalf("CreateFromDir: %v", err)
	}

	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()

	var sawServerFile bool
	for _, file := range zr.File {
		if strings.Contains(file.Name, "/cmd/create-templ-app/") {
			t.Errorf("module zip unexpectedly contains CLI file: %s", file.Name)
		}
		if strings.Contains(file.Name, "/cmd/server/main.go") {
			sawServerFile = true
		}
	}
	if !sawServerFile {
		t.Error("module zip is missing cmd/server/main.go — zip contents look wrong, not just over-filtered")
	}
}

// extractGitHEAD writes the git-tracked contents of HEAD in srcDir into
// dstDir, so tests operate on what's actually committed rather than
// whatever untracked build artifacts happen to be sitting on local disk.
func extractGitHEAD(t *testing.T, srcDir, dstDir string) {
	t.Helper()

	cmd := exec.Command("git", "archive", "--format=tar", "HEAD")
	cmd.Dir = srcDir
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	tr := tar.NewReader(pipe)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		target := filepath.Join(dstDir, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				t.Fatal(err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				t.Fatal(err)
			}
			out, err := os.Create(target)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				t.Fatal(err)
			}
			out.Close()
		}
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("git archive: %v", err)
	}
}
