package main

import (
	"embed"
	"path"
	"text/template"
)

//go:embed templates/*.template
var templates embed.FS

// Get a parsed template by name (e.g jtmp.ZFSinit)
func GetTemplate(name string) (tmp *template.Template, err error) {
	return template.ParseFS(templates, path.Join("templates", name))
}
