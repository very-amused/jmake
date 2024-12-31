package template

import (
	"bytes"
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
		} else {
			switch d := data.(type) {
			// FIXME
			case []any:
				var buf bytes.Buffer
				for _, obj := range d[:len(d)-1] {
					if err = tmp.Execute(&buf, obj); err != nil {
						log.Println(err)
					}
					tmp, err = template.New(templateName).Parse(buf.String())
				}
				if err = tmp.Execute(file, d[len(d)-1]); err != nil {
					log.Println(err)
				}

			default:
				if err = tmp.Execute(file, data); err != nil {
					log.Println(err)
				}
			}

			file.Close()
		}
	}
}
