package main

import (
	"bufio"
	"os"
	"path"

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

// #region legacy

func (z *ZFSconfig) makeTemplates() (err error) {
	if z.Dataset == "" {
		return nil
	}

	var (
		initScript    *bufio.Writer
		statusScript  *bufio.Writer
		destroyScript *bufio.Writer
	)

	if file, err := os.Create(path.Join(jtmp.AutoTemplateDir, jtmp.ZFSinit)); err != nil {
		return err
	} else {
		defer file.Close()
		initScript = bufio.NewWriter(file)
		defer initScript.Flush()
	}
	if file, err := os.Create(path.Join(jtmp.AutoTemplateDir, jtmp.ZFSstatus)); err != nil {
		return err
	} else {
		defer file.Close()
		statusScript = bufio.NewWriter(file)
		defer statusScript.Flush()
	}
	if file, err := os.Create(path.Join(jtmp.AutoTemplateDir, jtmp.ZFSdestroy)); err != nil {
		return err
	} else {
		defer file.Close()
		destroyScript = bufio.NewWriter(file)
		defer destroyScript.Flush()
	}

	initScript.WriteString("# Initialize jail ZFS datasets\n")

	statusScript.WriteString("# View status for jail ZFS datasets\n")

	destroyScript.WriteString("# Permanently destroy jail ZFS datasets (DANGER ZONE)\n")
	destroyScript.WriteString("# This will only work if all live jail datasets (in {{.Dataset}}/containers) have been deleted\n")

	// Create root dataset
	initScript.WriteString("{{if .Mountpoint}}\n")
	jtmp.WriteCommand(initScript, "zfs create -o mountpoint={{.Mountpoint}} {{.Dataset}}", true)
	initScript.WriteString("{{end}}\n")

	// Create child datasets
	jtmp.WriteCommand(initScript, "zfs create {{.Dataset}}/media", true)      // Compressed FreeBSD release images
	jtmp.WriteCommand(initScript, "zfs create {{.Dataset}}/templates", true)  // FreeBSD userland templates used to create thin jails via snapshot + clone
	jtmp.WriteCommand(initScript, "zfs create {{.Dataset}}/containers", true) // Live jail containers

	// Get status
	jtmp.WriteCommand(statusScript, "zfs list -r {{.Dataset}}", false)

	// Destroy datasets
	jtmp.WriteCommand(destroyScript, "zfs destroy {{.Dataset}}/containers", true) // Destroy containers non-recursively (only works if no live jail datasets exist)
	jtmp.WriteCommand(destroyScript, "[ $? = 0 ] || exit", false)                 // EXIT if the previous action failed (indicating live jail datasets)

	jtmp.WriteCommand(destroyScript, "zfs destroy -r {{.Dataset}}/templates", true) // Destroy templates (including snapshots)
	jtmp.WriteCommand(destroyScript, "zfs destroy -r {{.Dataset}}/media", true)     // Destroy downloaded install media
	jtmp.WriteCommand(destroyScript, "zfs destroy {{.Dataset}}", true)              // Finally, destroy the root dataset

	return nil
}

func (z *ZFSconfig) execTemplates() {
	jtmp.ExecAutoTemplates(z, jtmp.ZFSinit, jtmp.ZFSstatus, jtmp.ZFSdestroy)
}

// #endregion
