package main

import (
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
	//JailDeploy = "jail-deploy.sh.template"
)

// Get a template's output filename
func Output(template string) string {
	return strings.TrimSuffix(path.Base(template), ".template")
}
