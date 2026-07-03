package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func writeManifest(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestResolverURL_HashedMatch(t *testing.T) {
	path := writeManifest(t, `{"main.css": "main-abc123.css"}`)

	r, err := NewResolver(path, "/static")
	if err != nil {
		t.Fatalf("NewResolver: %v", err)
	}

	if got, want := r.URL("main.css"), "/static/main-abc123.css"; got != want {
		t.Errorf("URL(\"main.css\") = %q, want %q", got, want)
	}
}

func TestResolverURL_FallsBackForUnknownName(t *testing.T) {
	path := writeManifest(t, `{}`)

	r, err := NewResolver(path, "/static")
	if err != nil {
		t.Fatalf("NewResolver: %v", err)
	}

	if got, want := r.URL("missing.css"), "/static/missing.css"; got != want {
		t.Errorf("URL(\"missing.css\") = %q, want %q", got, want)
	}
}

func TestResolverURL_TrimsTrailingSlashFromBaseURL(t *testing.T) {
	path := writeManifest(t, `{"main.css": "main-abc123.css"}`)

	r, err := NewResolver(path, "/static/")
	if err != nil {
		t.Fatalf("NewResolver: %v", err)
	}

	if got, want := r.URL("main.css"), "/static/main-abc123.css"; got != want {
		t.Errorf("URL(\"main.css\") = %q, want %q", got, want)
	}
}

func TestNewResolver_MissingFile(t *testing.T) {
	_, err := NewResolver(filepath.Join(t.TempDir(), "nope.json"), "/static")
	if err == nil {
		t.Fatal("NewResolver: want error for missing manifest file, got nil")
	}
}

func TestNewResolver_InvalidJSON(t *testing.T) {
	path := writeManifest(t, `not json`)

	_, err := NewResolver(path, "/static")
	if err == nil {
		t.Fatal("NewResolver: want error for invalid JSON, got nil")
	}
}
