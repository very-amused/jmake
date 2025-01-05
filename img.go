package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/very-amused/jmake/jtmp"
)

type ImgConfig struct {
	Release  string // FreeBSD release string (e.g 14.2-RELEASE) to deploy from
	Snapshot string // A snapshot name to deploy from. Set to `base` if no image customization is desired
}

func (_ *ImgConfig) makeTemplates(c *Config) (err error) {
	if c.Img.Release == "" {
		return nil
	} else if c.ZFS == nil || c.ZFS.Dataset == "" || c.ZFS.Mountpoint == "" {
		return nil
	}

	var (
		initScript   *bufio.Writer
		removeScript *bufio.Writer
	)

	if file, err := os.Create(path.Join(jtmp.AutoTemplateDir, jtmp.ImgInit)); err != nil {
		return err
	} else {
		defer file.Close()
		initScript = bufio.NewWriter(file)
		defer initScript.Flush()
	}
	if file, err := os.Create(path.Join(jtmp.AutoTemplateDir, jtmp.ImgRemove)); err != nil {
		return err
	} else {
		defer file.Close()
		removeScript = bufio.NewWriter(file)
		defer removeScript.Flush()
	}

	// Absolute path to compressed imgTar.txz
	const imgTar = "{{.ZFS.Mountpoint}}/media/{{.Img.Release}}-base.txz"
	// Absolute path to extracted base.txz template
	const tmp = "{{.ZFS.Mountpoint}}/templates/{{.Img.Release}}"
	// ZFS dataset name for extracted base.txz template
	const tmpDataset = "{{.ZFS.Dataset}}/templates/{{.Img.Release}}"

	// #region img-init

	initScript.WriteString("# Initialize a FreeBSD {{.Img.Release}} image for customization and jail deployment\n\n")

	// Create image dataset
	jtmp.WriteCommand(initScript, fmt.Sprintf("zfs create %s", tmpDataset), true)

	// Download compressed image
	// TODO: These options should be configurable
	const mirror = "https://download.freebsd.org/ftp"
	const arch = "amd64/amd64"
	url := strings.Join([]string{mirror, arch, "{{.Img.Release}}", "base.txz"}, "/")
	jtmp.WriteCommand(initScript, fmt.Sprintf("fetch %s -o %s", url, imgTar), true)

	// Extract compressed image to dataset
	jtmp.WriteCommand(initScript, fmt.Sprintf("tar -xf %s -C %s --unlink", imgTar, tmp), true)

	// Configure and update template
	jtmp.WriteCommand(initScript, fmt.Sprintf("cp /etc/resolv.conf %s/etc/resolv.conf", tmp), true)
	jtmp.WriteCommand(initScript, fmt.Sprintf("cp /etc/localtime %s/etc/localtime", tmp), true)
	jtmp.WriteCommand(initScript, fmt.Sprintf("freebsd-update -b %s/ fetch install", tmp), true)

	// Create vanilla image snapshot
	jtmp.WriteCommand(initScript, fmt.Sprintf("zfs snapshot %s@vanilla", tmpDataset), true)

	// #endregion img-init

	// #region img-remove

	removeScript.WriteString("# Remove the downloaded FreeBSD {{.Img.Release}} image (requires that no jails depend depend on this image or its dataset)\n\n")

	// Remove dataset
	jtmp.WriteCommand(removeScript, fmt.Sprintf("zfs destroy -r %s", tmpDataset), true)
	jtmp.WriteCommand(removeScript, "[ $? = 0 ] || exit", false) // Exit if the previous action failed

	// Remove compressed image
	jtmp.WriteCommand(removeScript, fmt.Sprintf("rm -f %s", imgTar), true)

	// #endregion img-remove

	return nil
}

func (img *ImgConfig) execTemplates(c *Config) {
	jtmp.ExecAutoTemplates(c, jtmp.ImgInit, jtmp.ImgRemove)
}
