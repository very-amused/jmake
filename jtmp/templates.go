package jtmp

import (
	"path"
	"strings"
)

// Template names
const (
	ZFSinit    = "zfs-init.sh.template"
	ZFSstatus  = "zfs-status.sh.template"
	ZFSdestroy = "zfs-destroy.sh.template"

	ImgInit   = "img-init.sh.template"
	ImgRemove = "img-remove.sh.template"

	BridgeRC = "bridge.rc.conf.template"
)

// Get a template's output filename
func Output(template string) string {
	return strings.TrimSuffix(path.Base(template), ".template")
}
