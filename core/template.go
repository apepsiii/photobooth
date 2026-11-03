package core

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin/render"
)

// PageData is the standard data payload passed to every HTML template.
// Pages that need extra data can populate the Data field with a custom struct.
type PageData struct {
	Title     string
	UserName  string
	Flash     string
	FlashType string
	Path      string
	Data      interface{}
}

// multiRenderer implements gin's render.HTMLRender interface.
// Each page gets its own *template.Template set so that the {{block "content"}}
// in the master layout resolves to the correct page content without conflicts
// between pages that share the same block name.
type multiRenderer struct {
	templates map[string]*template.Template
}

// Instance returns a render.Render for the named template.
func (r *multiRenderer) Instance(name string, data interface{}) render.Render {
	tmpl, ok := r.templates[name]
	if !ok {
		panic(fmt.Sprintf("template not found: %s", name))
	}
	return render.HTML{
		Template: tmpl,
		Name:     name,
		Data:     data,
	}
}

// LoadTemplates builds a multiRenderer by parsing the master layout together
// with each individual page file. Every page gets an isolated template set so
// that block definitions ("content") never collide across pages.
func LoadTemplates(templateDir string) *multiRenderer {
	mr := &multiRenderer{templates: make(map[string]*template.Template)}

	layoutDir := filepath.Join(templateDir, "layouts")
	pagesDir := filepath.Join(templateDir, "pages")

	// Collect layout file paths.
	layouts, err := filepath.Glob(filepath.Join(layoutDir, "*.html"))
	if err != nil {
		panic("failed to glob layouts: " + err.Error())
	}

	pages, err := filepath.Glob(filepath.Join(pagesDir, "*.html"))
	if err != nil {
		panic("failed to glob pages: " + err.Error())
	}

	for _, pagePath := range pages {
		name := strings.TrimSuffix(filepath.Base(pagePath), ".html")

		// Start a fresh template set for this page.
		tmpl := template.New("")

		// Parse all layout files (defines "master", "auth", etc.).
		for _, l := range layouts {
			raw, err := os.ReadFile(l)
			if err != nil {
				panic("failed to read layout " + l + ": " + err.Error())
			}
			if _, err := tmpl.Parse(string(raw)); err != nil {
				panic("failed to parse layout " + l + ": " + err.Error())
			}
		}

		// Parse the page file (defines the page name + "content" block).
		raw, err := os.ReadFile(pagePath)
		if err != nil {
			panic("failed to read page " + pagePath + ": " + err.Error())
		}
		if _, err := tmpl.Parse(string(raw)); err != nil {
			panic("failed to parse page " + pagePath + ": " + err.Error())
		}

		mr.templates[name] = tmpl
	}

	return mr
}
