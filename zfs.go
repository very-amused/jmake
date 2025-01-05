package main

import (
	"github.com/very-amused/jmake/jtmp"
)

// ZFS dataset configuration for creating thin jail images/containers.
type ZFSconfig struct {
	Dataset    string // Jail root dataset
	Mountpoint string // Root jail mountpoint

	ContextChecks
}

func (z *ZFSconfig) Generate(_ *Config) (errs []error) {
	if z.Dataset == "" {
		return nil
	}

	if z.Mountpoint != "" {
		errs = append(errs, jtmp.ExecTemplates(z, jtmp.ZFSinit)...)
	}
	errs = append(errs, jtmp.ExecTemplates(z, jtmp.ZFSstatus, jtmp.ZFSdestroy)...)
	return errs
}
