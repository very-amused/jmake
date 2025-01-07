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
	Host   *HostConfig

	bridgeKeys []string
	jailKeys   []string // ordered jail config keys, not what's on the warden's belt

	ContextChecks
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
func (c *Config) parseKeyOrder() (e error) {
	file, e := os.Open("jmake.toml")
	if e != nil {
		return e
	}
	defer file.Close()

	parser := tomlu.Parser{}
	content, e := io.ReadAll(file)
	if e != nil {
		return e
	}
	parser.Reset(content)

	var table []string
	for parser.NextExpression() {
		exp := parser.Expression()

		var keyParts []string // Dot separated list of key parts
		switch exp.Kind {
		case tomlu.Table:
			table = splitTomlKey(exp.Key())
			keyParts = table
		case tomlu.KeyValue:
			parts := splitTomlKey(exp.Key())
			keyParts = append(table, parts...)
		default:
			continue
		}

		if len(keyParts) == 4 {
			switch keyParts[0] {
			case "jail":
				jail := keyParts[1]
				ip := keyParts[3]
				if c.Jail == nil || (*c.Jail)[jail] == nil {
					break
				}
				ipKeys := &(*c.Jail)[jail].ipKeys
				*ipKeys = append(*ipKeys, ip)
			}
			continue
		}
		if len(keyParts) == 2 {
			switch keyParts[0] {
			case "bridge":
				c.bridgeKeys = append(c.bridgeKeys, keyParts[1])
			case "jail":
				c.jailKeys = append(c.jailKeys, keyParts[1])
			}
		}
	}

	return nil
}

// Split a toml key iterator into a list of keys
// i.e my.key.here -> {"my", "key", "here"}
func splitTomlKey(it tomlu.Iterator) (keys []string) {
	for it.Next() {
		node := it.Node()
		keys = append(keys, string(node.Data))
	}
	return keys
}
