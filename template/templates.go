package template

import (
	"os"
)

// Root template directory
const TemplateDir = "templates"

// Template names
const ZFSinit = "zfs-init.sh.template"
const ZFSstatus = "zfs-status.sh.template"
const ZFSdestroy = "zfs-destroy.sh.template"

const ImgInit = "img-init.sh.template"
const ImgRemove = "img-remove.sh.template"

const BridgeRC = "bridge-rc.conf.template"

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
