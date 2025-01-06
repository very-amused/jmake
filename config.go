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

// Log newErrs and return append(errs, newErrs...)
func logErrs(newErrs, errs []error) []error {
	for _, err := range newErrs {
		log.Println(err)
	}
	return append(errs, newErrs...)
}

// Generate - Generate output files using jmake.toml
func (c *Config) Generate() {
	var errs []error
	if c.ZFS != nil {
		errs = logErrs(c.ZFS.Generate(c), errs)
	}
	if c.Img != nil {
		errs = logErrs(c.Img.Generate(c), errs)
	}
	if c.Bridge != nil {
		errs = logErrs(c.Bridge.Generate(c), errs)
	}

	if len(errs) > 0 {
		errNoun := "errors"
		if len(errs) == 1 {
			errNoun = "error"
		}
		log.Printf("jmake encountered %d %s while generating output files\n", len(errs), errNoun)
	}
}

// Parse jmake.toml
func ParseConfig() (c *Config, e error) {
	c = new(Config)
	_, e = toml.DecodeFile("jmake.toml", &c)
	return c, e
}
