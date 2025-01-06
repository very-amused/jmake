package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

// A full jmake.toml config
type Config struct {
	ZFS    *ZFSconfig
	Img    *ImgConfig
	Bridge BridgeConfigs
	Jail   map[string]*JailConfig

	ContextChecks
}

// A config section capable of template gen
type ConfigSection interface {
	makeTemplates(c *Config) error
	execTemplates(c *Config)
}

// ExecTemplates - Execute config templates previous generated using MakeTemplates
func (c *Config) ExecTemplates() {
	if c.ZFS != nil {
		c.ZFS.Generate(c)
	}
	if c.Img != nil {
		if errs := c.Img.Generate(c); len(errs) > 0 {
			for _, err := range errs {
				log.Println(err)
			}
		}
	}
	if c.Bridge != nil {
		c.Bridge.Generate(c)
	}
}

// Parse jmake.toml
func ParseConfig() (c *Config, e error) {
	c = new(Config)
	_, e = toml.DecodeFile("jmake.toml", &c)
	return c, e
}
