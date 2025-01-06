package main

import (
	"embed"
	"os"
	"path"
	"text/template"
)

//go:embed templates/*.template
var templates embed.FS

// Get a parsed template by name (e.g jtmp.ZFSinit)
func GetTemplate(name string) (tmp *template.Template, err error) {
	return template.ParseFS(templates, path.Join("templates", name))
}

// Get and execute a template using the provided data, writing to Output(name)
func ExecTemplates(data any, names ...string) (errs []error) {
	for _, name := range names {
		// Create output file
		file, err := os.Create(Output(name))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// Ensure all files get closed
		defer file.Close()

		// Load and execute template
		tmp, err := GetTemplate(name)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		err = tmp.Execute(file, data)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		file.Close()
	}

	return errs
}
