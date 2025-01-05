package main

import (
	"github.com/very-amused/jmake/jtmp"
)

// ZFS dataset configuration for creating thin jail images/containers.
type ZFSconfig struct {
	Dataset    string // Jail root dataset
	Mountpoint string // Root jail mountpoint
}

func (_ *ZFSconfig) Generate(c *Config) (errs []error) {
	if c.ZFS.Dataset == "" {
		return nil
	}

	if c.ZFS.Mountpoint != "" {
		errs = append(errs, jtmp.ExecTemplates(c, jtmp.ZFSinit)...)
	}
	errs = append(errs, jtmp.ExecTemplates(c, jtmp.ZFSstatus, jtmp.ZFSdestroy)...)
	return errs
}
