package views

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func testRenderer() Renderer {
	return NewRenderer(func(name string) string { return "/static/" + name })
}

func TestIndexRendersLayoutAndSections(t *testing.T) {
	r := testRenderer()

	var buf bytes.Buffer
	if err := r.Index().Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}

	html := buf.String()
	for _, want := range []string{
		"<title>Home</title>",
		"/static/main.css",
		"/static/main.js",
		"create-templ-app",
		"htmx (server round-trip)",
		"Alpine (client-side)",
		"Clicked 0 times (server)",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("Index() output missing %q\noutput:\n%s", want, html)
		}
	}
}

func TestCounterRendersCount(t *testing.T) {
	r := testRenderer()

	var buf bytes.Buffer
	if err := r.Counter(42).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Clicked 42 times (server)") {
		t.Errorf("Counter(42) output missing expected count text:\n%s", html)
	}
	if !strings.Contains(html, `hx-post="/api/count"`) {
		t.Errorf("Counter(42) output missing hx-post attribute:\n%s", html)
	}
}
