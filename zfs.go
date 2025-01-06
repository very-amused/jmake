package main

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
		errs = append(errs, ExecTemplates(z, ZFSinit)...)
	}
	errs = append(errs, ExecTemplates(z, ZFSstatus, ZFSdestroy)...)
	return errs
}
