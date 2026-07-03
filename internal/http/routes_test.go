package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/codypotter/create-templ-app/internal/assets"
	"github.com/gin-gonic/gin"
)

// newTestRouter builds a router wired up like the real server, backed by a
// temp dir standing in for the esbuild dist directory. extraFiles are
// written into that dir before the manifest is loaded, for tests that need
// an actual static asset to serve.
func newTestRouter(t *testing.T, extraFiles map[string]string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dir := t.TempDir()
	for name, content := range extraFiles {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	manifest := map[string]string{"main.css": "main-abc123.css", "main.js": "main-def456.js"}
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(manifestPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	resolver, err := assets.NewResolver(manifestPath, "/static")
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	Register(r, resolver, dir)
	return r
}

func TestIndexRoute(t *testing.T) {
	r := newTestRouter(t, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	if !strings.Contains(body, "create-templ-app") {
		t.Errorf("body missing expected heading:\n%s", body)
	}
	if !strings.Contains(body, "/static/main-abc123.css") {
		t.Errorf("body missing resolved asset URL:\n%s", body)
	}
}

var countRe = regexp.MustCompile(`Clicked (\d+) times`)

func extractCount(t *testing.T, body string) int64 {
	t.Helper()
	m := countRe.FindStringSubmatch(body)
	if m == nil {
		t.Fatalf("could not find count in body: %s", body)
	}
	n, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		t.Fatalf("parse count %q: %v", m[1], err)
	}
	return n
}

func TestCounterRouteIncrementsAcrossRequests(t *testing.T) {
	r := newTestRouter(t, nil)

	post := func() string {
		req := httptest.NewRequest(http.MethodPost, "/api/count", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
		return w.Body.String()
	}

	first := extractCount(t, post())
	second := extractCount(t, post())

	if second != first+1 {
		t.Errorf("count went from %d to %d, want a delta of 1", first, second)
	}
}

func TestStaticAssetServing(t *testing.T) {
	r := newTestRouter(t, map[string]string{"main-abc123.css": "body{}"})

	req := httptest.NewRequest(http.MethodGet, "/static/main-abc123.css", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if got, want := w.Body.String(), "body{}"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
	if got, want := w.Header().Get("Cache-Control"), "public, max-age=31536000, immutable"; got != want {
		t.Errorf("Cache-Control = %q, want %q", got, want)
	}
}
