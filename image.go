package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	jtmp "github.com/very-amused/jmake/template"
)

type ImageConfig struct {
	Release  string // FreeBSD release string (e.g 14.2-RELEASE) to deploy from
	Snapshot string // A snapshot name to deploy from. Set to `base` if no image customization is desired
}

func (img *ImageConfig) makeTemplates(z *ZFSconfig) (err error) {
	if img.Release == "" || z.Dataset == "" || z.Mountpoint == "" {
		return nil
	}

	var (
		initScript *bufio.Writer
		//removeScript *bufio.Writer
	)

	if file, err := os.Create(path.Join(jtmp.TemplateDir, jtmp.ImageInit)); err != nil {
		return err
	} else {
		defer file.Close()
		initScript = bufio.NewWriter(file)
		defer initScript.Flush()
	}

	initScript.WriteString("# Initialize a FreeBSD {{.Image.Release}} image for customization and jail deployment\n\n")

	// Absolute path to compressed base.txz
	const base = "{{.ZFS.Mountpoint}}/media/FreeBSD-{{.Image.Release}}-base.txz"
	// Absolute path to extracted base.txz template
	const tmp = "{{.ZFS.Mountpoint}}/templates/FreeBSD-{{.Image.Release}}"
	// ZFS dataset name for extracted base.txz template
	const tmpDataset = "{{.ZFS.Dataset}}/templates/FreeBSD-{{.Image.Release}}"

	initScript.WriteString("# Create image dataset\n")
	jtmp.WriteCommand(initScript, fmt.Sprintf("zfs create %s", tmpDataset), true)

	initScript.WriteString("# Download base.txz\n")
	// TODO: These options should be configurable
	const mirror = "https://download.freebsd.org/ftp"
	const arch = "amd64/amd64"
	url := strings.Join([]string{mirror, arch, img.Release, "base.txz"}, "/")
	jtmp.WriteCommand(initScript, fmt.Sprintf("fetch %s -o %s", url, base), true)

	initScript.WriteString("# Extract base.txz\n")
	jtmp.WriteCommand(initScript, fmt.Sprintf("tar -xf %s -C %s --unlink", base, tmp), true)

	initScript.WriteString("# Prepare and update template\n")
	jtmp.WriteCommand(initScript, fmt.Sprintf("cp /etc/resolv.conf %s/etc/resolv.conf", tmp), true)
	jtmp.WriteCommand(initScript, fmt.Sprintf("cp /etc/localtime %s/etc/localtime", tmp), true)
	jtmp.WriteCommand(initScript, fmt.Sprintf("freebsd-update -b %s/ fetch install", tmp), true)

	initScript.WriteString("# Create vanilla image snapshot (useful to undo future modifications)\n")
	jtmp.WriteCommand(initScript, fmt.Sprintf("zfs snapshot %s@vanilla", tmpDataset), true)

	return nil
}

func (img *ImageConfig) execTemplates(c *Config) {
	jtmp.ExecTemplates(c, jtmp.ImageInit, jtmp.ImageRemove)
}
