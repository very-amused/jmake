package main

import (
	"net/url"
	"path"
)

type ImgConfig struct {
	Release  string // FreeBSD release string (e.g 14.2-RELEASE) to deploy from
	Snapshot string // A snapshot name to deploy from. Set to `base` if no image customization is desired
	Arch     string // FreeBSD architecture string (default: "amd64/amd64")
	Mirror   string // FreeBSD download mirror (default: https://download.freebsd.org/ftp/releases/)

	ContextChecks

	zfs *ZFSconfig // ptr to ZFS config
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
	if img == nil {
		return nil
	}
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
		img.Mirror = "https://download.freebsd.org/ftp/releases/"
	}

	errs = append(errs, ExecTemplates(img, ImgInit, ImgStatus, ImgRemove)...)
	return errs
}
