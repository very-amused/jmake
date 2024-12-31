package main

import "github.com/BurntSushi/toml"

// A full jmake.toml config
type Config struct {
	ZFS      *ZFSconfig
	Programs string //
}

func (c *Config) makeTemplates() (errs []error) {
	errs = make([]error, 0)
	if c.ZFS != nil {
		if err := c.ZFS.makeTemplates(); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
func (c *Config) execTemplates() {
	if c.ZFS != nil {
		c.ZFS.execTemplates()
	}
}

// Parse jmake.toml
func ParseConfig() (c *Config) {
	c = new(Config)
	toml.DecodeFile("jmake.toml", &c)
	return c
}
