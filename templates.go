package main

import (
	"os"
	"path"
	"strings"
)

// Template names
const (
	ZFSinit    = "zfs-init.sh.template"
	ZFSstatus  = "zfs-status.sh.template"
	ZFSdestroy = "zfs-destroy.sh.template"

	ImgInit      = "img-init.sh.template"
	ImgStatus    = "img-status.sh.template"
	ImgRemove    = "img-remove.sh.template"
	ImgBootstrap = "img-bootstrap.sh.template"

	BridgeRC = "bridge.rc.conf.template"

	JailInit      = "jail-init.sh.template"
	JailBootstrap = "jail-bootstrap.sh.template"
	JailConf      = "jail.conf.template"
)

// Get a template's output filename
func Output(template string) string {
	return strings.TrimSuffix(path.Base(template), ".template")
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
			}
			file.Close()
		}
	}

	return errs
}
