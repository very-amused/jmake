package template

import (
	"os"
)

// A config section that generates template files and then executes them with its values
type TemplateGen interface {
	makeTemplates() error
	execTemplates() // Execute generated templates. Usually simply calls ExecTemplates
}

// Root template directory
const TemplateDir = "templates"

// Template names
const ZFSinit = "zfs-init.sh.template"
const ZFSstatus = "zfs-status.sh.template"
const ZFSdestroy = "zfs-destroy.sh.template"

const ImageInit = "img-init.sh.template"
const ImageRemove = "img-remove.sh.template"

// Create and clear templates dir
func CreateTemplateDir() (e error) {
	if e = os.RemoveAll(TemplateDir); e != nil {
		return e
	}
	if e = os.MkdirAll(TemplateDir, 0o755); e != nil {
		return e
	}

	return nil
}
