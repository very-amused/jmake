package main

import (
	"embed"
	"os"
	"path"
	"strings"
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
		// Calculate output file perms
		var perm os.FileMode = 0o644
		filename := Output(name)
		if path.Ext(filename) == ".sh" {
			perm += 0o111 // Add +x to executable files
		}
		defer os.Chmod(filename, perm) // Set file perms after writing (needed to create executable files)

		// Create output file
		file, err := os.OpenFile(Output(name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
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

// Get and execute templates using the provided slices for data and suffix labels
func ExecMultiTemplates[D any](data []D, labels []string, names ...string) (errs []error) {
	if len(data) != len(labels) {
		return errs
	}

	for _, name := range names {
		// Load template
		tmp, err := GetTemplate(name)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for i, label := range labels {

			// Add label to output filename
			var output strings.Builder
			nameParts := strings.Split(Output(name), ".")
			output.WriteString(nameParts[0])
			output.WriteRune('-')
			output.WriteString(label)
			output.WriteRune('.')
			output.WriteString(strings.Join(nameParts[1:], "."))

			// Open output file
			file, err := os.Create(output.String())
			if err != nil {
				errs = append(errs, err)
			}
			defer file.Close()

			// Execute template
			err = tmp.Execute(file, data[i])
			if err != nil {
				errs = append(errs, err)
				continue
			}
			file.Close()
		}
	}

	return errs
}
