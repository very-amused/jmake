package template

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

// Execute templates using text/template
func ExecTemplates(data any, templateNames ...string) {
	for _, templateName := range templateNames {
		tmp, err := template.New(templateName).ParseFiles(path.Join(TemplateDir, templateName))
		outfile := strings.TrimSuffix(templateName, ".template")
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// Clear out of date output files
				os.Remove(outfile)
				continue // Skip missing templates
			}
			log.Println(err)
		}

		file, err := os.Create(outfile)
		if err != nil {
			log.Println(err)
			return
		}
		if err = tmp.Execute(file, data); err != nil {
			log.Println(err)
		}
	}
}

// Execute templates using text/template, appending to previous output in destination files
func ExecAppend(data any, templateNames ...string) {
	for _, templateName := range templateNames {
		tmp, err := template.New(templateName).ParseFiles(path.Join(TemplateDir, templateName))
		outfile := strings.TrimSuffix(templateName, ".template")
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// Clear out of date output files
				os.Remove(outfile)
				continue // Skip missing templates
			}
			log.Println(err)
		}

		file, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			log.Println(err)
			return
		}
		if err = tmp.Execute(file, data); err != nil {
			log.Println(err)
		}
	}
}
