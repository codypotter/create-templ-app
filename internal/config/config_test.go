package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("ASSETS_DIST_PATH", "")
	t.Setenv("ASSET_BASE_URL", "")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.AssetsDistPath != "internal/assets/dist" {
		t.Errorf("AssetsDistPath = %q, want %q", cfg.AssetsDistPath, "internal/assets/dist")
	}
	if cfg.AssetBaseURL != "/static" {
		t.Errorf("AssetBaseURL = %q, want %q", cfg.AssetBaseURL, "/static")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("ASSETS_DIST_PATH", "custom/dist")
	t.Setenv("ASSET_BASE_URL", "https://cdn.example.com/assets/")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.AssetsDistPath != "custom/dist" {
		t.Errorf("AssetsDistPath = %q, want %q", cfg.AssetsDistPath, "custom/dist")
	}
	if cfg.AssetBaseURL != "https://cdn.example.com/assets" {
		t.Errorf("AssetBaseURL = %q, want trailing slash trimmed, got %q", cfg.AssetBaseURL, cfg.AssetBaseURL)
	}
}
