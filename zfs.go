package main

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/very-amused/jmake/template"
)

/* Reference: FreeBSD Handbook 4th Edition - Chapter 17. Jails and Containers */

// ZFS dataset configuration for creating thin jail images/containers.
type ZFSconfig struct {
	Dataset    string // Jail root dataset
	Mountpoint string // Root jail mountpoint
}

// ZFS child dataset paths (obtain via ZFSconfig.Datasets())
type ZFSDatasets struct {
	Media      string // Stores compressed vanilla FreeBSD images
	Templates  string // Stores jail template userlands + snapshots
	Containers string // Stores jail
}

func (z *ZFSconfig) makeTemplates() (err error) {
	if z.Dataset == "" {
		return nil
	}

	var (
		initScript *bufio.Writer
	)

	if file, err := os.Create(path.Join(template.TemplateDir, template.ZFSinit)); err != nil {
		return err
	} else {
		defer file.Close()
		initScript = bufio.NewWriter(file)
		defer initScript.Flush()
	}

	/* ref Handbook ch. 17.3.3 */
	// Create root dataset
	if z.Mountpoint != "" {
		cmd := fmt.Sprintf("zfs create -o mountpoint=%s %s", z.Mountpoint, z.Dataset)
		template.WriteCommand(initScript, cmd, true)
	}
	// Create child datasets
	datasets := z.Datasets()
	template.WriteCommand(initScript, fmt.Sprintf("zfs create %s", datasets.Media), true)
	template.WriteCommand(initScript, fmt.Sprintf("zfs create %s", datasets.Templates), true)
	template.WriteCommand(initScript, fmt.Sprintf("zfs create %s", datasets.Containers), true)

	return nil
}

// Return a ZFSDatasets struct with paths to ZFS child datasets. Returns nil if z.Dataset == ""
func (z *ZFSconfig) Datasets() (datasets *ZFSDatasets) {
	if z.Dataset == "" {
		return nil
	}

	datasets = new(ZFSDatasets)
	datasets.Media = path.Join(z.Dataset, "media")
	datasets.Templates = path.Join(z.Dataset, "templates")
	datasets.Containers = path.Join(z.Dataset, "containers")
	return datasets
}
