package main

import (
	"bufio"
	"os"
	"path"

	"github.com/very-amused/jmake/template"
)

// ZFS dataset configuration for creating thin jail images/containers.
type ZFSconfig struct {
	Dataset    string // Jail root dataset
	Mountpoint string // Root jail mountpoint
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

	initScript.WriteString("# Reference: FreeBSD Handbook 4th Edition - Chapter 17. Jails and Containers\n")
	// Create root dataset
	initScript.WriteString("{{if .Mountpoint}}\n")
	template.WriteCommand(initScript, "zfs create -o mountpoint={{.Mountpoint}} {{.Dataset}}", true)
	initScript.WriteString("{{end}}\n")

	// Create child datasets
	template.WriteCommand(initScript, "zfs create {{.Dataset}}/media", true)      // Compressed FreeBSD release images
	template.WriteCommand(initScript, "zfs create {{.Dataset}}/templates", true)  // FreeBSD userland templates used to create thin jails via snapshot + clone
	template.WriteCommand(initScript, "zfs create {{.Dataset}}/containers", true) // Live jail containers

	return nil
}
