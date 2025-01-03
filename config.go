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

func (c *Config) makeTemplates() (errs []error) {
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
func (c *Config) execTemplates() {
	if c.ZFS != nil {
		c.ZFS.execTemplates()

		if c.Img != nil {
			c.Img.execTemplates(c)
		}
	}
	if len(c.Bridge) > 0 {
		outfile := strings.TrimSuffix(jtmp.BridgeRC, ".template")
		os.Remove(outfile)
		for i := range c.Bridge {
			c.Bridge[i].execTemplates(c)
		}
		content, _ := os.ReadFile(outfile)
		if file, err := os.Create(outfile); err == nil {
			rc := bufio.NewWriter(file)

			var ifaces []string
			for i := range c.Bridge {
				ifaces = append(ifaces, c.Bridge[i].Name)
			}
			jtmp.WriteRc(rc, "cloned_interfaces", strings.Join(ifaces, " "))
			rc.Write(content)
			rc.Flush()
			file.Close()
		}
	}
}

// Parse jmake.toml
func ParseConfig() (c *Config) {
	c = new(Config)
	toml.DecodeFile("jmake.toml", &c)
	return c
}
