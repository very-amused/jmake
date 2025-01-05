package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/very-amused/jmake/jtmp"
)

type ImgConfig struct {
	Release  string // FreeBSD release string (e.g 14.2-RELEASE) to deploy from
	Snapshot string // A snapshot name to deploy from. Set to `base` if no image customization is desired
	Arch     string // FreeBSD architecture string (default: "amd64/amd64")
	Mirror   string // FreeBSD download mirror (default: https://download.freebsd.org/ftp/)

	ContextChecks

	zfs *ZFSconfig `toml:"-"` // ptr to ZFS config
}

// Path to image root folder (where the image's base.txz is extracted)
func (img *ImgConfig) Path() string {
	return path.Join(img.zfs.Mountpoint, "templates", img.Release)
}

// Path to compressed image base tarball
func (img *ImgConfig) Tar() string {
	return path.Join(img.zfs.Mountpoint, "media", img.Release+"-base.txz")
}

// Base tarball download URL
func (img *ImgConfig) TarURL() string {
	tarURL, _ := url.JoinPath(img.Mirror, img.Arch, img.Release, "base.txz")
	return tarURL
}

// Name of image ZFS dataset
func (img *ImgConfig) Dataset() string {
	return path.Join(img.zfs.Dataset, "templates", img.Release)
}

func (img *ImgConfig) Generate(c *Config) (errs []error) {
	if img.Release == "" {
		return nil
	}
	if c.ZFS == nil || c.ZFS.Mountpoint == "" || c.ZFS.Dataset == "" {
		return nil
	}
	img.zfs = c.ZFS

	// Set defaults (TODO: better default handling)
	if img.Arch == "" {
		img.Arch = "amd64/amd64"
	}
	if img.Mirror == "" {
		img.Mirror = "https://download.freebsd.org/ftp/"
	}

	errs = append(errs, jtmp.ExecTemplates(img, jtmp.ImgInit)...)
	return errs
}

// #region legacy

func (_ *ImgConfig) makeTemplates(c *Config) (err error) {
	if c.Img.Release == "" {
		return nil
	} else if c.ZFS == nil || c.ZFS.Dataset == "" || c.ZFS.Mountpoint == "" {
		return nil
	}

	var (
		removeScript *bufio.Writer
	)

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
	jtmp.ExecAutoTemplates(c, jtmp.ImgRemove)
}

// #endregion
