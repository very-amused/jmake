package main

import (
	"log"
	"os"

	"github.com/very-amused/jmake/jtmp"
)

func main() {
	// Output to stderr
	log.SetOutput(os.Stderr)

	// Init template dir
	var err error
	if err = jtmp.CreateTemplateDir(); err != nil {
		log.Fatalln("Failed to initialize template dir:", err)
	}

	// Parse jmake.conf
	jmake, err := ParseConfig()
	if err != nil {
		log.Fatalln("Failed to parse jmake.toml:", err)
	}

	// Write templates
	jmake.MakeTemplates()
	// Execute templates
	jmake.ExecTemplates()
}
