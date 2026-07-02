package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/codypotter/create-templ-app/internal/assets"
	"github.com/codypotter/create-templ-app/internal/config"
	"github.com/codypotter/create-templ-app/internal/http"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	manifestPath := filepath.Join(cfg.AssetsDistPath, "manifest.json")
	resolver, err := assets.NewResolver(manifestPath, cfg.AssetBaseURL)
	if err != nil {
		log.Fatalf("asset resolver: %v", err)
	}

	r := gin.Default()
	http.Register(r, resolver, cfg.AssetsDistPath)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
