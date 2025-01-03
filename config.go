package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	jtmp "github.com/very-amused/jmake/template"
)

// A full jmake.toml config
type Config struct {
	ZFS    *ZFSconfig
	Img    *ImgConfig
	Bridge []BridgeConfig
}

// A config section capable of template gen
type ConfigSection interface {
	makeTemplates(c *Config) error
	execTemplates(c *Config)
}

// MakeTemplates - Make templates which can be executed with the loaded config
func (c *Config) MakeTemplates() (errs []error) {
	errs = make([]error, 0)
	if c.ZFS != nil {
		if err := c.ZFS.makeTemplates(); err != nil {
			errs = append(errs, err)
		}
	}
	if c.Img != nil {
		if err := c.Img.makeTemplates(c); err != nil {
			errs = append(errs, err)
		}
	}
	for i := range c.Bridge {
		if err := c.Bridge[i].makeTemplates(c); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// ExecTemplates - Execute config templates previous generated using MakeTemplates
func (c *Config) ExecTemplates() {
	if c.ZFS != nil {
		c.ZFS.execTemplates()

		if c.Img != nil {
			c.Img.execTemplates(c)
		}
	}
	if len(c.Bridge) > 0 {
		// Write cloned_interfaces rc
		outfile := strings.TrimSuffix(jtmp.BridgeRC, ".template")
		if file, err := os.Create(outfile); err == nil {
			rc := bufio.NewWriter(file)

			var ifaces []string
			for i := range c.Bridge {
				ifaces = append(ifaces, c.Bridge[i].Name)
			}
			jtmp.WriteRc(rc, "cloned_interfaces", strings.Join(ifaces, " "))
			rc.Flush()
			file.Close()
		}
		// Write ifconfig for bridge interfaces
		for i := range c.Bridge {
			c.Bridge[i].execTemplates(c)
		}
	}
}

// Parse jmake.toml
func ParseConfig() (c *Config) {
	c = new(Config)
	toml.DecodeFile("jmake.toml", &c)
	return c
}
