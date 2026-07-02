package assets

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Resolver maps logical asset names (e.g. "main.css") to their full URL,
// handling content-hashed filenames from the esbuild manifest.
type Resolver struct {
	baseURL  string
	manifest map[string]string
}

func NewResolver(manifestPath, baseURL string) (*Resolver, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read asset manifest: %w", err)
	}

	var manifest map[string]string
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse asset manifest: %w", err)
	}

	return &Resolver{
		baseURL:  strings.TrimRight(baseURL, "/"),
		manifest: manifest,
	}, nil
}

// URL returns the full URL for a logical asset name.
// e.g. URL("main.css") -> "/static/main-a1b2c3.css" or "https://assets.com/main-a1b2c3.css"
func (r *Resolver) URL(name string) string {
	hashed, ok := r.manifest[name]
	if !ok {
		return r.baseURL + "/" + name
	}
	return r.baseURL + "/" + hashed
}
