package main

import (
	"io"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	tomlu "github.com/pelletier/go-toml/v2/unstable"
)

// A full jmake.toml config
type Config struct {
	ZFS    *ZFSconfig
	Img    *ImgConfig
	Bridge *BridgeConfigs
	Jail   *JailConfigs

	Host *HostConfig

	ContextChecks

	Keys *ConfigKeys `toml:"-"`
}

// Ordered jmake.toml keys (needed b/c Go maps are unordered)
type ConfigKeys struct {
	Bridge []string
	Jail   []string
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
	file, e := os.Open("jmake.toml")
	if e != nil {
		return c, e
	}
	defer file.Close()

	dec := toml.NewDecoder(file)
	e = dec.Decode(&c)
	return c, e
}

// Parse jmake.toml key order
func ParseConfigKeys() (k *ConfigKeys, e error) {
	k = new(ConfigKeys)
	file, e := os.Open("jmake.toml")
	if e != nil {
		return k, e
	}
	defer file.Close()

	parser := tomlu.Parser{}
	content, e := io.ReadAll(file)
	if e != nil {
		return k, e
	}
	parser.Reset(content)

	for parser.NextExpression() {
		exp := parser.Expression()

		if exp.Kind != tomlu.Table {
			continue
		}

		key := exp.Key()
		keyParts := keyAsStrings(key)
		if len(keyParts) != 2 {
			continue
		}
		switch keyParts[0] {
		case "bridge":
			k.Bridge = append(k.Bridge, keyParts[1])
		case "jail":
			k.Jail = append(k.Jail, keyParts[1])
		}
	}

	return k, e
}

// helper to transfor a key iterator to a slice of strings
// credit: pelletier [https://github.com/pelletier/go-toml/discussions/801#discussioncomment-7083586]
func keyAsStrings(it tomlu.Iterator) []string {
	var parts []string
	for it.Next() {
		n := it.Node()
		parts = append(parts, string(n.Data))
	}
	return parts
}
