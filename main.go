package main

import (
	"log"
	"os"

	"github.com/very-amused/jmake/template"
)

func main() {
	// Output to stderr
	log.SetOutput(os.Stderr)

	// Init template dir
	var err error
	if err = template.CreateTemplateDir(); err != nil {
		log.Fatalln("Failed to initialize template dir:", err)
	}

	// Parse jmake.conf
	jmake := ParseConfig()

	// Write templates
	jmake.MakeTemplates()
	// Execute templates
	jmake.ExecTemplates()
}
