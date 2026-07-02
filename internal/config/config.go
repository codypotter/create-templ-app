package config

import (
	"os"
	"strings"
)

type Config struct {
	// Port the HTTP server listens on.
	Port string

	// AssetsDistPath is the local filesystem path to the esbuild output directory.
	// Used to serve /static/* and to find manifest.json.
	AssetsDistPath string

	// AssetBaseURL is the base URL prepended to resolved asset filenames.
	// In dev this is "/static". In prod it points at the CDN.
	AssetBaseURL string
}

func Load() Config {
	return Config{
		Port:           envOr("PORT", "8080"),
		AssetsDistPath: envOr("ASSETS_DIST_PATH", "internal/assets/dist"),
		AssetBaseURL:   strings.TrimRight(envOr("ASSET_BASE_URL", "/static"), "/"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
