package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

// A full jmake.toml config
type Config struct {
	ZFS    *ZFSconfig
	Img    *ImgConfig
	Bridge *BridgeConfigs
	Jail   *JailConfigs

	Host *HostConfig

	ContextChecks
}

// Log newErrs and return append(errs, newErrs...)
func logErrs(newErrs, errs []error) []error {
	for _, err := range newErrs {
		log.Println(err)
	}
	return append(errs, newErrs...)
}

// ConfigSection - A config section capable of generating output files
type ConfigSection interface {
	Generate(*Config) []error
}

// Generate - Generate output files using jmake.toml
func (c *Config) Generate() {
	// Generate output for all config sections supporting .Generate(c)
	var errs []error
	configSections := []ConfigSection{c.ZFS, c.Img, c.Bridge, c.Jail}
	for _, section := range configSections {
		if section == nil {
			continue
		}
		errs = logErrs(section.Generate(c), errs)
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
