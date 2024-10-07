package openapi

import (
	"embed"
	"net/http"

	"github.com/mvrilo/go-redoc"
)

// Config allows serving an OpenAPI
// documentation from within a http server
type Config struct {
	DocsPath    string
	SpecPath    string
	SpecFile    string
	SpecFS      *embed.FS
	Title       string
	Description string
}

func (c Config) HandlerFn() http.HandlerFunc {
	doc := redoc.Redoc{
		DocsPath:    c.DocsPath,
		SpecPath:    c.SpecPath,
		SpecFile:    c.SpecFile,
		SpecFS:      c.SpecFS,
		Title:       c.Title,
		Description: c.Description,
	}

	return doc.Handler()
}
