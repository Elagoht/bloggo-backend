package static

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type StaticModule struct {
	distFS fs.FS
}

func NewModule(distFS fs.FS) StaticModule {
	return StaticModule{
		distFS: distFS,
	}
}

func (module StaticModule) RegisterModule(router *chi.Mux) {
	// Serve static assets (CSS, JS, images, fonts, etc.)
	router.Get("/assets/*", module.ServeAssets)

	// Serve other static files at root (favicon, manifest, etc.)
	router.Get("/favicon.ico", module.ServeStaticFile)
	router.Get("/manifest.json", module.ServeStaticFile)
	router.Get("/robots.txt", module.ServeStaticFile)

	// Serve index.html for all other routes (SPA support)
	router.NotFound(module.ServeSPA)
}

func (module StaticModule) ServeAssets(w http.ResponseWriter, r *http.Request) {
	// Serve files from the dist filesystem
	http.FileServer(http.FS(module.distFS)).ServeHTTP(w, r)
}

func (module StaticModule) ServeStaticFile(w http.ResponseWriter, r *http.Request) {
	// Serve individual static files from root
	http.FileServer(http.FS(module.distFS)).ServeHTTP(w, r)
}

func (module StaticModule) ServeSPA(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Don't serve SPA for API routes or internal routes
	if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/internal/") {
		http.NotFound(w, r)
		return
	}

	// Check if the requested path is a static file
	if strings.HasPrefix(path, "/assets/") {
		http.FileServer(http.FS(module.distFS)).ServeHTTP(w, r)
		return
	}

	// For all other routes, serve index.html (SPA routing)
	indexPath := "index.html"
	content, err := fs.ReadFile(module.distFS, indexPath)
	if err != nil {
		http.Error(w, "Frontend not found", http.StatusNotFound)
		return
	}

	// Set proper content type for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	w.Write(content)
}
