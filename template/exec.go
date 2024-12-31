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
		tmp, err := template.ParseFiles(path.Join(TemplateDir, templateName))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue // Skip missing templates
			}
			log.Println(err)
		}

		outfile := strings.TrimSuffix(templateName, ".template")
		file, err := os.Create(outfile)
		if err != nil {
			log.Println(err)
		} else {
			if err := tmp.Execute(file, data); err != nil {
				log.Println(err)
			}
			file.Close()
		}
	}
}
